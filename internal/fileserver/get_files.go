package fileserver

import (
	"context"

	"github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
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
	// 1) найти пользователя по email
	userIDRaw := ctx.Value("UserID")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	u, err := s.Q.ListFilesByUser(ctx, uid) // напиши этот sqlc-метод
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	out := make([]*filev1.FileItem, 0, len(u))
	for _, r := range u {
		out = append(out, &filev1.FileItem{
			Filename:  r.Filename,
			Size:      r.SizeBytes,
			CreatedAt: timestamppb.New(r.CreatedAt.UTC()),
		})
	}

	return &filev1.GetFilesResponse{
		Files: out,
	}, nil

}
