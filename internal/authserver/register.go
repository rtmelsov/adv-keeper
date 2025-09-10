package authserver

import (
	"context"
	"fmt"

	"errors"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	db "github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/helpers"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	commonv1.UnimplementedAuthServiceServer
	Q *db.Queries
}

func New(q *db.Queries) *Service { return &Service{Q: q} }

func (s *Service) Register(ctx context.Context, in *commonv1.RegisterRequest) (*commonv1.RegisterResponse, error) {
	email := strings.TrimSpace(in.GetEmail())
	pass := in.GetPassword()
	if email == "" || pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	_, err := s.Q.GetUserByEmail(ctx, email) // напиши этот sqlc-метод
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("пользователь уже существует"))
	}

	// простая защита: хэшируем пароль (для MVP достаточно bcrypt)
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash error")
	}

	// внутри Register(...)
	arg := db.RegisterParams{
		Email:   email,
		PwdPhc:  string(hash), // bcrypt/argon2 — как у тебя
		E2eePub: nil,          // или []byte{} / из запроса
	}

	ID, err := s.Q.Register(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "db error")
	}

	// 3) сгенерить токены
	access, exp, err := helpers.NewAccessJWT(ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "jwt")
	}

	expiresAt := timestamppb.New(exp)

	return &commonv1.RegisterResponse{
		UserId: ID.String(),
		Tokens: &commonv1.TokenPair{
			AccessToken: access,
			ExpiresAt:   expiresAt,
		},
	}, nil
}
