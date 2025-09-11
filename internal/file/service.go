// Package file
package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	db "github.com/rtmelsov/adv-keeper/internal/db"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
)

type FileServer struct {
	filev1.UnimplementedFileServiceServer
	Q         *db.Queries
	uploadDir string
}

func New(q *db.Queries) *FileServer {
	envs, err := helpers.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return &FileServer{Q: q, uploadDir: envs.FilesDir}
}

func (s *FileServer) Upload(stream filev1.FileService_UploadServer) error {
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
	info := first.GetInfo()
	if info == nil {
		return fmt.Errorf("first message must be FileInfo")
	}

	filename = filepath.Base(info.Filename)
	if filename == "" {
		filename = fmt.Sprintf("upload-%d.bin", time.Now().UnixNano())
	}

	tmpPath := filepath.Join(s.uploadDir, filename+".part")
	finalPath := filepath.Join(s.uploadDir, filename)

	// гарантируем папку
	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return err
	}
	outFile, err = os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer func() {
		outFile.Close()
		// Если контекст отменён — удалим частичный файл
		if stream.Context().Err() != nil {
			_ = os.Remove(tmpPath)
		}
	}()

	// 2) принимаем куски
	for {
		// Проверим отмену клиента/дедлайн
		if err := stream.Context().Err(); err != nil {
			return err
		}

		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ch := msg.GetChunk()
		if ch == nil {
			return fmt.Errorf("unexpected message type: want FileChunk")
		}

		n, err := outFile.Write(ch.Content)
		if err != nil {
			return err
		}
		written += int64(n)
		_, _ = hasher.Write(ch.Content) // считаем sha256 «на лету»
	}

	// 3) закрываем и переименовываем atomic-стилем
	if err := outFile.Sync(); err != nil {
		return err
	}
	if err := outFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return err
	}

	sum := hex.EncodeToString(hasher.Sum(nil))
	log.Printf("Upload done: %s, bytes=%d, sha256=%s, took=%s",
		finalPath, written, sum, time.Since(startTime))

	return stream.SendAndClose(&filev1.UploadResponse{
		StoredAs:      finalPath,
		BytesReceived: written,
		Sha256:        sum,
	})
}
