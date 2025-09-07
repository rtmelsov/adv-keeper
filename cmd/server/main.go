package main

import (
	"database/sql"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/auth"
	db "github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/file"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	// получаем урл бд

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к Postgres
	dbx, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	q := db.New(dbx)
	s := grpc.NewServer()
	commonv1.RegisterAuthServiceServer(s, auth.New(q))
	filev1.RegisterFileServiceServer(s, file.New(q))

	reflection.Register(s)

	log.Println("gRPC listening on", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
