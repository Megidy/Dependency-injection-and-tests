package main

import (
	"database/sql"
	"log"

	"github.com/API/cmd/api"
	"github.com/API/db"
)

func main() {
	db, err := db.NewMySqlStorage("root:password@tcp(mysql:3306)/api")
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	server := api.NewApiServer(":8080", db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)

	}
	log.Println("DB: successfully connected!")
}
