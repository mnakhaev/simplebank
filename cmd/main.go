package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mnakhaev/simplebank/api"
	"github.com/mnakhaev/simplebank/config"
	db "github.com/mnakhaev/simplebank/db/sqlc"
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
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
