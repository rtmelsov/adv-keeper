package akclient

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
	"github.com/rtmelsov/adv-keeper/internal/models"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
)

const chunkSize = 1 << 20 // 1 MiB — безопасно ниже 4 MiB лимита на сообщение

func UploadFile(path string, prog chan<- models.Prog) {
	defer close(prog)
	ctx, err := middleware.AddAuthData()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	f, err := os.Open(path)
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	total := stat.Size()
	conn, err := grpc.NewClient(
		helpers.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			// Можно поднять лимит на отправку одного сообщения при желании
			grpc.MaxCallSendMsgSize(8<<20),
		),
	)

	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}
	defer conn.Close()

	client := filev1.NewFileServiceClient(conn)

	stream, err := client.Upload(ctx)
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	// 1) отправляем мета-инфу (первое сообщение)
	err = stream.Send(&filev1.UploadRequest{
		Payload: &filev1.UploadRequest_Info{
			Info: &filev1.FileInfo{
				Filename: filepath.Base(path),
				Size:     stat.Size(),
			},
		},
	})
	if err != nil {
		prog <- models.Prog{Err: err}
		return
	}

	// 2) шлём файл кусками
	buf := make([]byte, chunkSize)
	var offset int64
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			chunk := &filev1.FileChunk{
				Content: buf[:n],
				Offset:  offset,
			}
			if err := stream.Send(&filev1.UploadRequest{
				Payload: &filev1.UploadRequest_Chunk{Chunk: chunk},
			}); err != nil {
				prog <- models.Prog{Err: err}
				return
			}
			offset += int64(n)
			select {
			case prog <- models.Prog{Done: offset, Total: total}:
			default: // не блокируем UI, если буфер заполнен
			}
		}
		if readErr == io.EOF {
			break // <— обязательно выходим!
		}
		if readErr != nil {
			prog <- models.Prog{Err: err}
			return
		}
	}

	// 3) закрываем отправку и получаем ответ
	_, err = stream.CloseAndRecv()
	if err != nil {
		if st, ok := status.FromError(err); ok {
			prog <- models.Prog{Err: fmt.Errorf("upload failed: %s: %s", st.Code(), st.Message())}
			return
		}
		prog <- models.Prog{Err: fmt.Errorf("upload failed: %w", err)}
		return
	}
}
