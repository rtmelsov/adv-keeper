package fileserver

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	db "github.com/rtmelsov/adv-keeper/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/charmbracelet/log"
	"github.com/rtmelsov/adv-keeper/internal/helpers"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) DeleteFile(ctx context.Context, req *filev1.DeleteFileRequest) (*emptypb.Empty, error) {
	navs, err := helpers.LoadConfig()
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return nil, err
	}

	uid, err := helpers.UserIDFromCtx(ctx)
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return nil, err
	}

	id, err := helpers.StringParseUUID(req.Fileid)
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return nil, err
	}

	u, err := s.Q.GetFileForUser(ctx, db.GetFileForUserParams{
		UserID: uid,
		ID:     id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db: %v", err)
	}

	// 1) удаляем запись и получаем path
	_, err = s.Q.DeleteFile(ctx, db.DeleteFileParams{
		ID:     id,
		UserID: uid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db: %v", err)
	}

	// 2) чуть-чуть безопасной валидации пути (чтобы не удалить лишнее)
	baseDir := navs.FilesDir // задай в конструкторе из ENV: /opt/adv-keeper/data
	absBase, _ := filepath.Abs(baseDir)
	absPath, _ := filepath.Abs(u.Path)
	if !strings.HasPrefix(absPath, absBase+string(os.PathSeparator)) && absPath != absBase {
		// если вдруг в БД путь «вне» каталога — не трогаем ФС
		// но запись уже удалили — это окей, просто залогируй
		// log.Warn("path outside base", "path", absPath)
		return &emptypb.Empty{}, nil
	}

	// 3) удаляем файл на ФС (если нет — не ошибка)
	_ = os.Remove(absPath)
	return &emptypb.Empty{}, nil
}
