// Package storepg
package storepg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	dbpkg "github.com/rtmelsov/adv-keeper/internal/db" // <- замени путь
)

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// RegisterWithDevice Обёртка вокруг sqlc-запроса RegisterWithDevice.
// Возвращает user_id и device_id строками.
func (r *Repo) RegisterWithDevice(ctx context.Context, args dbpkg.RegisterParams) (string, error) {
	ID, err := r.Q.Register(ctx, args)
	if err != nil {
		if isUnique(err) {
			return "", errors.New("конфликт по уникальному email")
			// конфликт по уникальному email — отдадим как есть, обработаем выше
		}
		return "", err
	}
	// не завязываемся на конкретный тип UUID — приводим к строке
	return toString(ID), nil
}

func toString(v any) string { return fmt.Sprint(v) }
