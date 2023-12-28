package main

import (
	"database/sql"
	"log"

	"github.com/TruongCongToan/simple_bank/api"
	db "github.com/TruongCongToan/simple_bank/db/sqlc"
	"github.com/TruongCongToan/simple_bank/util"

	_ "github.com/lib/pq" //without it code can not be able to talk to database
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal(" Can not load config")
		return
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not connect to database:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(" Can not create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can not start server:", err)
	}
}
