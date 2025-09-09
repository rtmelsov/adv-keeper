package middleware

import (
	"context"
	"google.golang.org/grpc/metadata"

	"github.com/rtmelsov/adv-keeper/internal/helpers"
)

func AddAuthData() (context.Context, error) {
	session, err := helpers.LoadSession()

	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"authorization": "Bearer " + session.AccessToken})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx, nil
}
