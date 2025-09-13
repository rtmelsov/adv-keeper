package akclient

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func DownloadFile(fileID string, prog chan<- models.Prog) {
	defer close(prog)

	ctx, err := middleware.AddAuthData()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	outDir, err := helpers.GetDownloadsDir()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	conn, err := grpc.DialContext(ctx, helpers.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}
	defer conn.Close()

	c := filev1.NewFileServiceClient(conn)
	stream, err := c.DownloadFile(ctx, &filev1.DownloadFileRequest{Fileid: fileID})
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

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

	filename := helpers.NextAvailableName(outDir, info.GetFilename())
	total := info.GetSize()

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
	var retErr error
	defer func() {
		_ = out.Close()
		if retErr != nil {
			_ = os.Remove(tmpPath) // удаляем именно .part при ошибке
		}
	}()

	h := sha256.New()
	var written int64
	var eofHex string

	for {
		if ctx.Err() != nil {
			retErr = ctx.Err()
			prog <- models.Prog{Err: retErr}
			return
		}

		msg, rerr := stream.Recv()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			if st, ok := status.FromError(rerr); ok {
				retErr = fmt.Errorf("recv: %s", st.Message())
			} else {
				retErr = fmt.Errorf("recv: %w", rerr)
			}
			prog <- models.Prog{Err: retErr}
			return
		}

		if e := msg.GetEof(); e != nil {
			eofHex = strings.ToLower(strings.TrimSpace(e.GetSha256Hex()))
			continue
		}
		if ch := msg.GetChunk(); ch != nil {
			n, werr := out.Write(ch.Content)
			if werr != nil {
				retErr = fmt.Errorf("write: %w", werr)
				prog <- models.Prog{Err: retErr}
				return
			}
			written += int64(n)
			select {
			case prog <- models.Prog{Done: written, Total: total}:
			default:
			}
			_, _ = h.Write(ch.Content)
			continue
		}

		retErr = errors.New("unexpected message: want chunk or eof")
		prog <- models.Prog{Err: retErr}
		return
	}

	gotHex := strings.ToLower(hex.EncodeToString(h.Sum(nil)))

	if eofHex == "" {
		retErr = errors.New("missing EOF sha256 from server")
		prog <- models.Prog{Err: retErr}
		return
	}
	if !strings.EqualFold(gotHex, eofHex) {
		retErr = fmt.Errorf("sha256 mismatch: got %s, want %s", gotHex, eofHex)
		prog <- models.Prog{Err: retErr}
		return
	}
	if total > 0 && written != total {
		retErr = fmt.Errorf("size mismatch: got %d, want %d", written, total)
		prog <- models.Prog{Err: retErr}
		return
	}

	if err := out.Sync(); err != nil {
		retErr = fmt.Errorf("sync: %w", err)
		prog <- models.Prog{Err: retErr}
		return
	}
	if err := out.Close(); err != nil {
		retErr = fmt.Errorf("close: %w", err)
		prog <- models.Prog{Err: retErr}
		return
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		retErr = fmt.Errorf("rename: %w", err)
		prog <- models.Prog{Err: retErr}
		return
	}
}
