// Package middleware
package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/charmbracelet/log"
)

const (
	headerAuthorize = "authorization"
)

var allowUnauth = map[string]struct{}{
	"/common.v1.AuthService/Register": {},
	"/common.v1.AuthService/Login":    {},
	"/common.v1.AuthService/Logout":   {},
}

func ServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if _, ok := allowUnauth[info.FullMethod]; ok {
		return handler(ctx, req) // пропускаем без проверки
	}

	vals := metadata.ValueFromIncomingContext(ctx, headerAuthorize)
	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "Request unauthenticated with ")
	}

	_, token, found := strings.Cut(vals[0], " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "Bad authorization string")
	}

	m, err := helpers.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "UserID", m.UserID)

	resp, err := handler(newCtx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StreamInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Info("get into stream interceptor")

	// Разрешённые методы без авторизации
	if _, ok := allowUnauth[info.FullMethod]; ok {
		return handler(srv, ss)
	}

	// Достаём токен
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		return status.Error(codes.Unauthenticated, "missing auth token")
	}

	_, token, found := strings.Cut(vals[0], " ")
	if !found {
		return status.Error(codes.Unauthenticated, "bad authorization string")
	}

	log.Info("access", token)
	claims, err := helpers.VerifyToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// Заворачиваем stream с новым контекстом
	newCtx := context.WithValue(ss.Context(), "UserID", claims.UserID)
	wrapped := &wrappedStream{ss, newCtx}

	return handler(srv, wrapped)
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
