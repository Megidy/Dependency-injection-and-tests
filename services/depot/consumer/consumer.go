package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/API/types"
	"github.com/IBM/sarama"
)

var wg sync.WaitGroup

const (
	ConsumerTopic string = "send_orders"
)

type KafkaConsumer struct {
	consumer   sarama.Consumer
	depotStore types.DepotStore
}

func NewConsumer(consumer sarama.Consumer, depotStore types.DepotStore) *KafkaConsumer {
	return &KafkaConsumer{
		consumer:   consumer,
		depotStore: depotStore,
	}
}

func (c *KafkaConsumer) ReceiveOrders() error {
	consumer, err := c.consumer.ConsumePartition(ConsumerTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	doneCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)

			case message := <-consumer.Messages():

				var order types.Order
				err := json.Unmarshal(message.Value, &order)
				if err != nil {
					log.Println(err)
				}
				log.Println("received order ", order.Product, "from user : ", order.UserID)
				wg.Add(1)
				go c.HandleOrder(order, &wg)

			case <-sigCh:
				log.Println("interrupt")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("exited goroutine")
	err = c.consumer.Close()
	if err != nil {
		return err
	}
	return nil
}
func (c *KafkaConsumer) HandleOrder(order types.Order, wg *sync.WaitGroup) {
	defer wg.Done()
	order.Status = "Processing"
	err := c.depotStore.UpdateOrderStatus(order)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < 50; i++ {
		fmt.Println("order of user: ", order.UserID, " is ready for: ", i+1, "%")
		time.Sleep(time.Second)
	}
	order.Status = "Ready to pickup"
	err = c.depotStore.UpdateOrderStatus(order)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("order of user ", order.UserID, " is ready!")

}
