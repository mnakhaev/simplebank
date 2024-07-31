package grpc_api

import (
	"testing"
	"time"

	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
	"github.com/mnakhaev/simplebank/util"
	"github.com/mnakhaev/simplebank/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, td worker.TaskDistributor) *Server {
	cfg := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(cfg, store, td)
	require.NoError(t, err)

	return server
}
