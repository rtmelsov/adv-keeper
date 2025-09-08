package middleware

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func AddAuthData(AccessToken string) context.Context {

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + AccessToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}
