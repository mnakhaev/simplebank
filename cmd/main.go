package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/mnakhaev/simplebank/api"
	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/grpc_api"
	"github.com/mnakhaev/simplebank/pb"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load cfg:", err)
	}

	dbConn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewSQLStore(dbConn)
	runGRPCServer(cfg, store)
}

func runGRPCServer(cfg config.Config, store db.Store) {
	server, err := grpc_api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	// reflection.Register(grpcServer) // allows gRPC client to explore which RPC are available on the server and how to call them.

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGinServer(cfg config.Config, store db.Store) {
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
