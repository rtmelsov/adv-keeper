package middleware

import (
	"context"

	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc/metadata"
)

func AddAuthData() (context.Context, error) {
	token := ""
	session, err := helpers.LoadSession()
	if err != nil {
		return nil, err
	}
	if session != nil {
		token = session.AccessToken
	}
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx, nil
}
