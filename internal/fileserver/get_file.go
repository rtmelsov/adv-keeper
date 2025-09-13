package fileserver

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/google/uuid"

	db "github.com/rtmelsov/adv-keeper/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/charmbracelet/log"
	"github.com/rtmelsov/adv-keeper/internal/helpers"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
)

const chunkSize = 1 << 20 // 1 MiB — безопасно ниже 4 MiB лимита на сообщение

func (s *Service) DownloadFile(DownloadFileRequest *filev1.DownloadFileRequest, stream filev1.FileService_DownloadFileServer) error {
	ctx := stream.Context()
	uid, err := helpers.UserIDFromCtx(ctx)
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return err
	}

	id, err := uuid.Parse(DownloadFileRequest.Fileid)
	if err != nil {
		log.Error("error while try to uuid.Parse", "error", err.Error())
		return err
	}

	u, err := s.Q.GetFileForUser(ctx, db.GetFileForUserParams{
		UserID: uid,
		ID:     id,
	})
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}
	path := u.Path
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	// 1) отправляем мета-инфу (первое сообщение)
	err = stream.Send(&filev1.DownloadFileResponse{
		Payload: &filev1.DownloadFileResponse_Info{
			Info: &filev1.FileInfo{
				Filename: u.Filename,
				Size:     stat.Size(),
			},
		},
	})
	if err != nil {
		return err
	}

	// 2) шлём файл кусками
	buf := make([]byte, chunkSize)
	h := sha256.New()
	var sent int64

	for {
		// отмена клиента/дедлайн
		if ctx.Err() != nil {
			return ctx.Err()
		}

		n, rerr := f.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			sent += int64(n)
			_, _ = h.Write(chunk)

			if err := stream.Send(&filev1.DownloadFileResponse{
				Payload: &filev1.DownloadFileResponse_Chunk{
					Chunk: &filev1.FileChunk{Content: chunk},
				},
			}); err != nil {
				return err
			}
		}
		if rerr == io.EOF {
			if err := stream.Send(&filev1.DownloadFileResponse{
				Payload: &filev1.DownloadFileResponse_Eof{
					Eof: &filev1.FileEof{Sha256Hex: hex.EncodeToString(h.Sum(nil))},
				},
			}); err != nil {
				return err
			}

			// 3) закрываем отправку и получаем ответ
			log.Info("download done",
				"file", u.Path,
				"bytes", u.SizeBytes,
				"sha256", hex.EncodeToString(h.Sum(nil)),
			)
			return nil
		}
		if rerr != nil {
			return status.Errorf(codes.Internal, "read: %v", rerr)
		}
	}

	return nil
}
