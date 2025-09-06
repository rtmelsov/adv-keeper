package helpers

import (
	"database/sql"

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
