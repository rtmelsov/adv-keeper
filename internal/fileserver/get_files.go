package fileserver

import (
	"context"

	"github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
)

type Service struct {
	filev1.UnimplementedFileServiceServer
	Q         *db.Queries
	UploadDir string
}

func New(q *db.Queries) *Service {
	conf, _ := helpers.LoadConfig()
	return &Service{Q: q, UploadDir: conf.FilesDir}
}

func (s *Service) GetFiles(ctx context.Context, GetFileRequest *filev1.GetFilesRequest) (*filev1.GetFilesResponse, error) {
	uid, err := helpers.UserIDFromCtx(ctx)
	if err != nil {
		log.Error("error while try to get client id", "error", err.Error())
		return nil, err
	}

	u, err := s.Q.ListFilesByUser(ctx, uid) // напиши этот sqlc-метод
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	out := make([]*filev1.FileItem, 0, len(u))
	for _, r := range u {
		out = append(out, &filev1.FileItem{
			Fileid:    r.ID.String(),
			Filename:  r.Filename,
			Size:      r.SizeBytes,
			CreatedAt: timestamppb.New(r.CreatedAt.UTC()),
		})
	}

	log.Info("files", "list", out)

	return &filev1.GetFilesResponse{
		Files: out,
	}, nil

}
