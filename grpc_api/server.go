package grpc_api

import (
	"fmt"

	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/pb"
	"github.com/mnakhaev/simplebank/token"
	"github.com/mnakhaev/simplebank/worker"
)

// Server serves gRPC requests for banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer // needed to receive gRPC calls before they are actually implemented
	config                           config.Config
	store                            db.Store
	tokenMaker                       token.Maker
	taskDistributor                  worker.TaskDistributor
}

// NewServer creates new gRPC server.
func NewServer(config config.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}
	return server, nil
}
