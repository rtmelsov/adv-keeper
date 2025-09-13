package helpers

import (
	"context"
	"database/sql"
	db "github.com/rtmelsov/adv-keeper/internal/db"
)

func withTx(ctx context.Context, dbx *sql.DB, fn func(q *db.Queries) error) error {
	tx, err := dbx.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // можно повысить до Serializable при гонках
	})
	if err != nil {
		return err
	}

	qtx := db.New(tx) // sqlc: создаём Queries, «подключённый» к транзакции

	if err := fn(qtx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
