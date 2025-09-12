// Package akclient
package akclient

import (
	"log"

	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
)

func GetProfile() (*commonv1.GetProfileResponse, error) {
	ctx, err := middleware.AddAuthData()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(helpers.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial %s: %v", helpers.Addr, err)
	}
	defer conn.Close()

	// 2) gRPC-клиент
	c := commonv1.NewAuthServiceClient(conn)

	return c.GetProfile(ctx, &emptypb.Empty{})
}
