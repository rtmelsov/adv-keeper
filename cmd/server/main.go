package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/auth"
	db "github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/file"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

func main() {
	// получаем урл бд
	dsn := "postgres://postgres@localhost:5432/dbname?sslmode=disable"

	fmt.Println("dsn", dsn)

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
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
