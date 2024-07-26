package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	_ "github.com/mnakhaev/simplebank/doc/statik" // needed to embed static files into the binary
	"github.com/mnakhaev/simplebank/mail"
	"github.com/mnakhaev/simplebank/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/mnakhaev/simplebank/api"
	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/grpc_api"
	"github.com/mnakhaev/simplebank/pb"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load cfg")
	}

	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // for pretty log printing
	}

	dbConn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(cfg.MigrationURL, cfg.DBSource)

	store := db.NewSQLStore(dbConn)

	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(cfg, redisOpt, store)

	go runGatewayServer(cfg, store, taskDistributor)
	runGRPCServer(cfg, store, taskDistributor)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration instance")
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}
	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config config.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	sender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, sender)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGRPCServer(cfg config.Config, store db.Store, td worker.TaskDistributor) {
	server, err := grpc_api.NewServer(cfg, store, td)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(grpc_api.GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer) // allows gRPC client to explore which RPC are available on the server and how to call them.

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(cfg config.Config, store db.Store, td worker.TaskDistributor) {
	server, err := grpc_api.NewServer(cfg, store, td)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	// use camel cases for gateway response instead of snake case
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux) // handle all endpoint by grpcMux

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := grpc_api.HTTPLogger(mux) // add middleware for HTTP requests logging.
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func runGinServer(cfg config.Config, store db.Store) {
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
