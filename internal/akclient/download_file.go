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
	"github.com/rtmelsov/adv-keeper/internal/models"
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
func DownloadFile(fileID string, prog chan<- models.Prog) {
	defer close(prog)
	ctx, err := middleware.AddAuthData()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}
	outDir := helpers.DownloadFilesDir
	conn, err := grpc.NewClient(helpers.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		prog <- models.Prog{Err: err}
		log.Fatalf("dial %s: %v", helpers.Addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := filev1.NewFileServiceClient(conn)

	stream, err := c.DownloadFile(ctx, &filev1.DownloadFileRequest{Fileid: fileID})
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	// 1) ждём FileInfo
	first, err := stream.Recv()
	if err != nil {
		prog <- models.Prog{Err: status.Errorf(codes.Internal, "read: %v", err)}
		return
	}
	info := first.GetInfo()
	if info == nil {
		prog <- models.Prog{Err: errors.New("first message must be FileInfo")}
		return
	}

	filename := safeBase(info.GetFilename())
	if filename == "" {
		filename = "file"
	}

	// 2) готовим пути/директории
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		prog <- models.Prog{Err: fmt.Errorf("mkdir: %w", err)}
		return
	}
	tmpPath := filepath.Join(outDir, filename+".part")
	finalPath := filepath.Join(outDir, filename)

	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		prog <- models.Prog{Err: fmt.Errorf("create: %w", err)}
		return
	}
	defer func() {
		out.Close()
		if err != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	stat, err := out.Stat()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}
	total := stat.Size()

	h := sha256.New()
	var written int64

	for {
		if ctx.Err() != nil {
			prog <- models.Prog{Err: ctx.Err()}
			return
		}
		msg, rerr := stream.Recv()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			// красиво покажем gRPC статус
			if st, ok := status.FromError(rerr); ok {
				prog <- models.Prog{Err: fmt.Errorf("recv: %s", st.Message())}
				return
			}
			prog <- models.Prog{Err: fmt.Errorf("recv: %w", rerr)}
			return
		}

		ch := msg.GetChunk()
		if ch == nil {
			prog <- models.Prog{Err: errors.New("unexpected message: want FileChunk")}
			return
		}

		n, werr := out.Write(ch.Content)
		if werr != nil {
			prog <- models.Prog{Err: fmt.Errorf("write: %w", werr)}
			return
		}
		written += int64(n)
		select {
		case prog <- models.Prog{Done: written, Total: total}:
		default: // не блокируем UI, если буфер заполнен
		}
		_, _ = h.Write(ch.Content)
	}

	// 4) fsync → close → rename
	if err := out.Sync(); err != nil {
		prog <- models.Prog{Err: fmt.Errorf("sync: %w", err)}
		return
	}
	if err := out.Close(); err != nil {
		prog <- models.Prog{Err: fmt.Errorf("close: %w", err)}
		return
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		prog <- models.Prog{Err: fmt.Errorf("rename: %w", err)}
		return
	}

	// 5) сверим размер (если сервер прислал size)
	if info.Size > 0 && written != info.Size {
		prog <- models.Prog{Err: fmt.Errorf("size mismatch: got %d, want %d", written, info.Size)}
		return
	}
}
