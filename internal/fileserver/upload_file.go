// Package fileserver
package fileserver

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	db "github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/helpers"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
)

func (s *Service) Upload(stream filev1.FileService_UploadServer) error {
	log.Info("try to get client id")

	ctx := stream.Context()
	uid, err := helpers.UserIDFromCtx(ctx)
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return err
	}
	log.Info("upload get request...")
	var (
		outFile   *os.File
		written   int64
		filename  string
		hasher    = sha256.New()
		startTime = time.Now()
	)

	// 1) ждём первое сообщение с метаданными
	first, err := stream.Recv()
	if err != nil {
		return err
	}

	log.Info("stream recv...")
	info := first.GetInfo()
	if info == nil {
		return fmt.Errorf("first message must be FileInfo")
	}

	log.Info("first get info...")

	filename = filepath.Base(info.Filename)
	if filename == "" {
		filename = fmt.Sprintf("upload-%d.bin", time.Now().UnixNano())
	}

	log.Info("file path - file name...", "path", filename)

	tmpPath := filepath.Join(s.UploadDir, filename+".part")
	finalPath := filepath.Join(s.UploadDir, filename)

	log.Info("start: mkdir...", "file name", filename)

	log.Info("mkdir all uplaod dir...", "path", s.UploadDir)
	outFile, err = os.Create(tmpPath)
	if err != nil {

		log.Error("OS CREATE", "error: ", err.Error())
		return err
	}

	log.Info("created dir...", s.UploadDir)
	defer func() {
		outFile.Close()
		// Если контекст отменён — удалим частичный файл
		if ctx.Err() != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	log.Info("start to get file info...")
	// 2) принимаем куски
	for {
		log.Info("uploading by piece...")
		// Проверим отмену клиента/дедлайн
		if err := ctx.Err(); err != nil {
			log.Error("Context", "Error", err.Error())
			return err
		}

		log.Info("1")

		msg, err := stream.Recv()

		log.Info("1.2")
		if err == io.EOF {
			log.Error("io.EOF: stream Recv", "Error", err.Error())
			break
		}

		log.Info("1.5")
		if err != nil {
			log.Error("stream Recv", "Error", err.Error())
			return err
		}

		log.Info("2")
		ch := msg.GetChunk()
		if ch == nil {
			log.Error("unexpected message type: want FileChunk")
			return fmt.Errorf("unexpected message type: want FileChunk")
		}

		n, err := outFile.Write(ch.Content)
		if err != nil {
			log.Error("outFile.Write", "Error", err.Error())
			return err
		}

		log.Info("3")
		written += int64(n)
		_, _ = hasher.Write(ch.Content) // считаем sha256 «на лету»
		log.Info("end to add that piece...")
	}

	log.Info("end: getting info...")
	// 3) закрываем и переименовываем atomic-стилем
	if err := outFile.Sync(); err != nil {
		return err
	}

	log.Info("out file close...")
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return err
	}
	log.Info("os rename...")

	sum := hex.EncodeToString(hasher.Sum(nil))
	log.Info("Upload done: %s, bytes=%d, sha256=%s, took=%s",
		finalPath, written, sum, time.Since(startTime))

	_, err = s.Q.AddFile(ctx, db.AddFileParams{
		UserID:    uid,
		Filename:  filename,
		Path:      finalPath,
		SizeBytes: info.GetSize(),
	})
	if err != nil {
		_ = os.Remove(finalPath) // best effort

		return status.Errorf(codes.Internal, "db insert failed: %v", err)
	}

	return stream.SendAndClose(&filev1.UploadResponse{
		StoredAs:      finalPath,
		BytesReceived: written,
		Sha256:        sum,
	})
}
