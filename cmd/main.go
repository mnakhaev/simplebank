package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mnakhaev/simplebank/db/util"

	"github.com/mnakhaev/simplebank/api"
	db "github.com/mnakhaev/simplebank/db/sqlc"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewSQLStore(dbConn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
