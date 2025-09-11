package helpers

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

func Nstr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{} // Valid=false => в БД NULL
	}
	return sql.NullString{String: s, Valid: true}
}

func Nuuid(s string) uuid.NullUUID {
	if s == "" {
		return uuid.NullUUID{} // Valid=false => в БД NULL (триггерит COALESCE)
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.NullUUID{} // можно и вернуть ошибку, если хочешь жёстче
	}
	return uuid.NullUUID{UUID: id, Valid: true}
}

// где-то в пакете auth
func UserIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value("UserID") // или строковый ключ, который кладёшь в интерсепторе
	s, ok := v.(string)
	if !ok || s == "" {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing user id")
	}
	return StringParseUUID(s)
}

func StringParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	return id, err
}
