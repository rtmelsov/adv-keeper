package akclient

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
)

func GetFiles() (*filev1.GetFilesResponse, error) {
	GetFilesRequest := &filev1.GetFilesRequest{
		Limit: 50,
	}

	ctxWithMeta, err := middleware.AddAuthData()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(helpers.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial %s: %v", helpers.Addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := filev1.NewFileServiceClient(conn)

	// 3) Вызываем RPC (клиент инициирует запрос)
	ctx, cancel := context.WithTimeout(ctxWithMeta, 40*time.Second)
	defer cancel()

	return c.GetFiles(ctx, GetFilesRequest)
}
