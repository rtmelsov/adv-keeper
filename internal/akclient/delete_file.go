package akclient

import (
	"fmt"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
)

func DeleteFile(fileID string) error {
	log.Info("delete file")
	if fileID == "" {
		return fmt.Errorf("empty fileID")
	}
	envs, err := helpers.LoadConfig()
	if err != nil {
		return err
	}

	// ctx с Authorization: Bearer <token>
	ctx, err := middleware.AddAuthData()
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(envs.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial %s: %w", envs.Addr, err)
	}
	defer conn.Close()

	c := filev1.NewFileServiceClient(conn)

	_, err = c.DeleteFile(ctx, &filev1.DeleteFileRequest{Fileid: fileID})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return fmt.Errorf("delete: %s", st.Message())
		}
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
