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
func (r *Repo) RegisterWithDevice(ctx context.Context, args dbpkg.RegisterWithDeviceParams) (string, string, error) {
	row, err := r.Q.RegisterWithDevice(ctx, args)
	if err != nil {
		if isUnique(err) {
			// конфликт по уникальному email — отдадим как есть, обработаем выше
		}
		return "", "", err
	}
	// не завязываемся на конкретный тип UUID — приводим к строке
	return toString(row.UserID), toString(row.DeviceID), nil
}

func toString(v any) string { return fmt.Sprint(v) }
