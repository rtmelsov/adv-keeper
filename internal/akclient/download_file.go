package akclient

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"

	"errors"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func safeBase(name string) string {
	// убираем директории и опасные символы
	base := filepath.Base(name)
	runes := make([]rune, 0, len(base))
	for _, r := range base {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '.', r == '-', r == '_', r == ' ':
			runes = append(runes, r)
		}
	}
	if len(runes) == 0 {
		return "file"
	}
	// защитим длину
	if len(runes) > 128 {
		runes = runes[:128]
	}
	return string(runes)
}

func DownloadFile(fileID string) (*filev1.GetFilesResponse, error) {
	envs, err := helpers.LoadConfig()

	if err != nil {
		return nil, err
	}
	ctx, err := middleware.AddAuthData()
	if err != nil {
		return nil, err
	}
	outDir := envs.DownloadFilesDir
	conn, err := grpc.NewClient(envs.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial %s: %v", envs.Addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := filev1.NewFileServiceClient(conn)

	stream, err := c.DownloadFile(ctx, &filev1.DownloadFileRequest{Fileid: fileID})
	if err != nil {
		return nil, err
	}

	// 1) ждём FileInfo
	first, err := stream.Recv()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "read: %v", err)
	}
	info := first.GetInfo()
	if info == nil {
		return nil, errors.New("first message must be FileInfo")
	}

	filename := safeBase(info.GetFilename())
	if filename == "" {
		filename = "file"
	}

	// 2) готовим пути/директории
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	tmpPath := filepath.Join(outDir, filename+".part")
	finalPath := filepath.Join(outDir, filename)

	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer func() {
		out.Close()
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	h := sha256.New()
	var written int64

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		msg, rerr := stream.Recv()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			// красиво покажем gRPC статус
			if st, ok := status.FromError(rerr); ok {
				return nil, fmt.Errorf("recv: %s", st.Message())
			}
			return nil, fmt.Errorf("recv: %w", rerr)
		}

		ch := msg.GetChunk()
		if ch == nil {
			return nil, errors.New("unexpected message: want FileChunk")
		}

		n, werr := out.Write(ch.Content)
		if werr != nil {
			return nil, fmt.Errorf("write: %w", werr)
		}
		written += int64(n)
		_, _ = h.Write(ch.Content)
	}

	// 4) fsync → close → rename
	if err := out.Sync(); err != nil {
		return nil, fmt.Errorf("sync: %w", err)
	}
	if err := out.Close(); err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return nil, fmt.Errorf("rename: %w", err)
	}

	// 5) сверим размер (если сервер прислал size)
	if info.Size > 0 && written != info.Size {
		return nil, fmt.Errorf("size mismatch: got %d, want %d", written, info.Size)
	}

	return nil, nil
}
