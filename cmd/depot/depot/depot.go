package depot

import (
	"database/sql"
	"log"

	"github.com/API/services/depot"
	"github.com/API/services/depot/consumer"
	"github.com/IBM/sarama"
)

type DepotServer struct {
	db       *sql.DB
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewDepotServer(db *sql.DB, producer sarama.SyncProducer, consumer sarama.Consumer) *DepotServer {
	return &DepotServer{
		db:       db,
		producer: producer,
		consumer: consumer,
	}
}

func (d *DepotServer) Run() error {

	log.Println("STARTED DEPOT SERVER")
	depotStore := depot.NewStore(d.db)
	kafkaConsumer := consumer.NewConsumer(d.consumer, depotStore)
	err := kafkaConsumer.ReceiveOrders()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
