package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/reflection"

	commonv1 "github.com/rtmelsov/adv-keeper/gen/go/proto/common/v1"
	filev1 "github.com/rtmelsov/adv-keeper/gen/go/proto/file/v1"
	"github.com/rtmelsov/adv-keeper/internal/authserver"
	db "github.com/rtmelsov/adv-keeper/internal/db"
	"github.com/rtmelsov/adv-keeper/internal/fileserver"
	"github.com/rtmelsov/adv-keeper/internal/helpers"
	"github.com/rtmelsov/adv-keeper/internal/middleware"
	"google.golang.org/grpc/keepalive"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

func main() {
	envs, err := helpers.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", envs.Addr)
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к Postgres
	dbx, err := sql.Open("pgx", envs.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	q := db.New(dbx)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.ServerInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     0,
			MaxConnectionAge:      0,
			MaxConnectionAgeGrace: 0,
			Time:                  2 * time.Minute,
			Timeout:               20 * time.Second,
		}),
	)
	commonv1.RegisterAuthServiceServer(s, authserver.New(q))
	filev1.RegisterFileServiceServer(s, fileserver.New(q))

	reflection.Register(s)

	log.Println("gRPC listening on", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
