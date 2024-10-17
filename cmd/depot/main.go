package main

import (
	"database/sql"
	"log"

	"github.com/API/cmd/depot/depot"
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

	consumerBrokers := []string{ConsumerPort}
	kafkaConsumer, err := kafka.ConnectConsumer(consumerBrokers)
	if err != nil {
		log.Fatal(err)
	}

	depotServer := depot.NewDepotServer(db, kafkaProducer, kafkaConsumer)
	log.Println("Starting DEPOT service")

	err = depotServer.Run()
	if err != nil {
		log.Fatal(err)
	}

}
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)

	}
	log.Println("DEPOT DB: successfully connected!")
}
