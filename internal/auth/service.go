// Package auth
package auth

import (
	"context"
	"errors"
	"strings"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	db "github.com/rtmelsov/adv-keeper/internal/db"

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

	// простая защита: хэшируем пароль (для MVP достаточно bcrypt)
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash error")
	}

	// внутри Register(...)
	arg := db.RegisterWithDeviceParams{
		Email:    email,
		PwdPhc:   string(hash), // bcrypt/argon2 — как у тебя
		E2eePub:  nil,          // или []byte{} / из запроса
		DeviceID: "",
	}

	row, err := s.Q.RegisterWithDevice(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "db error")
	}

	return &commonv1.RegisterResponse{
		UserId:   row.UserID.String(),
		DeviceId: row.DeviceID,
		Email:    arg.Email,
	}, nil
}

func (s *Service) Login(ctx context.Context, in *commonv1.LoginRequest) (*commonv1.LoginResponse, error) {
	email := strings.TrimSpace(in.GetEmail())
	pass := in.GetPassword()
	if email == "" || pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	// простая защита: хэшируем пароль (для MVP достаточно bcrypt)
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash error")
	}

	// внутри Register(...)
	arg := db.RegisterWithDeviceParams{
		Email:    email,
		PwdPhc:   string(hash), // bcrypt/argon2 — как у тебя
		E2eePub:  nil,          // или []byte{} / из запроса
		DeviceID: "",
	}

	row, err := s.Q.RegisterWithDevice(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "db error")
	}

	return &commonv1.LoginResponse{
		UserId:   row.UserID.String(),
		DeviceId: row.DeviceID,
		Email:    arg.Email,
	}, nil
}
