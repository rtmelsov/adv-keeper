// Package akclient
package akclient

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func Register(RegisterRequest *commonv1.RegisterRequest) (*commonv1.RegisterResponse, error) {
	// 1) Соединение с сервером
	addr := "127.0.0.1:8080" // или где у тебя слушает сервер
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial %s: %v", addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := commonv1.NewAuthServiceClient(conn)

	// 3) Вызываем RPC (клиент инициирует запрос)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	return c.Register(ctx, RegisterRequest)
}
