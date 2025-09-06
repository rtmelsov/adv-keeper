package storepg

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для database/sql

	dbpkg "github.com/rtmelsov/adv-keeper/internal/db" // <- замени на свой модульный путь
)

type Repo struct {
	DB *sql.DB
	Q  *dbpkg.Queries
}

func Open(dsn string) (*Repo, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &Repo{DB: db, Q: dbpkg.New(db)}, nil
}

func (r *Repo) Close() error { return r.DB.Close() }
