// Package akclient
package akclient

import (
	"context"
	"log"
	"time"

	"errors"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func Register(RegisterRequest *commonv1.RegisterRequest) (*commonv1.RegisterResponse, error) {
	envs, err := helpers.LoadConfig()
	if err != nil {
		return nil, errors.New("не получилось распарсить переменные окуржения")
	}
	conn, err := grpc.NewClient(envs.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial %s: %v", envs.Addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := commonv1.NewAuthServiceClient(conn)

	// 3) Вызываем RPC (клиент инициирует запрос)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	return c.Register(ctx, RegisterRequest)
}
