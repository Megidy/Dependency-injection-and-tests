package api

import (
	"database/sql"
	"log"

	"github.com/API/services/depot/producer"
	"github.com/API/services/order"
	"github.com/API/services/product"
	"github.com/API/services/user"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort  string = "kafka:9092"
	ProducerTopic string = "send_orders"
)

type ApiServer struct {
	addr     string
	db       *sql.DB
	producer sarama.SyncProducer
}

func NewApiServer(addr string, db *sql.DB, producer sarama.SyncProducer) *ApiServer {
	return &ApiServer{
		addr:     addr,
		db:       db,
		producer: producer,
	}
}

func (s *ApiServer) Run() error {
	router := gin.Default()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(router)

	newProducer := producer.NewProducer(s.producer)
	orderStore := order.NewStore(s.db)
	orderHandler := order.NewHandler(orderStore, productStore, userStore, newProducer)
	orderHandler.RegisterRoutes(router)

	log.Println("started server on 8080 ")
	return router.Run()

}
