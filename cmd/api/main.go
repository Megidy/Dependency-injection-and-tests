package main

import (
	"database/sql"
	"log"

	"github.com/API/cmd/api/api"
	"github.com/API/db"
	"github.com/API/kafka"
)

const (
	ConsumerPort  string = "kafka:9092"
	ConsumerTopic string = "receive_orders"
	ProducerPort  string = "kafka:9092"
	ProducerTopic string = "send_orders"
)

func main() {
	db, err := db.NewMySqlStorage("root:password@tcp(mysql:3306)/api")
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	producerBrokers := []string{ProducerPort}
	kafkaProducer, err := kafka.ConnectProducer(producerBrokers)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(":8080", db, kafkaProducer)
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
	log.Println("API DB: successfully connected!")
}
