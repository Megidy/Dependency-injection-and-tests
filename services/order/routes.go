package order

import (
	"encoding/json"
	"net/http"

	"github.com/API/services/auth"
	"github.com/API/types"
	"github.com/API/utils"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort  string = "kafka:9092"
	ProducerTopic string = "send_orders"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
	producer     types.Producer
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore, producer types.Producer) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
		producer:     producer,
	}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {

	router.POST("/order", auth.NewHandler(h.userStore).WithJWT, h.CreateOrder)
}
func (h *Handler) CreateOrder(c *gin.Context) {
	user, _ := c.Get("user")
	var payload types.CreateOrderPayload
	var order types.Order
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		utils.HandleError(c, err, "failed to read body", http.StatusBadRequest)
		return
	}
	product, err := h.productStore.GetProductById(payload.Id)
	if err != nil {
		utils.HandleError(c, err, "failed to get sinle product", http.StatusBadRequest)
		return
	}
	if product.Id == 0 {
		utils.HandleError(c, err, "no product with this id", http.StatusBadRequest)
		return
	}
	if payload.Quantity > product.Quantity {
		utils.HandleError(c, err, "not avaible in this amount", http.StatusBadRequest)
		return
	}
	order = types.Order{
		UserID:  user.(*types.User).Id,
		Product: product,
		Status:  "Pending",
	}

	err = h.store.CreateOrder(order)
	if err != nil {
		utils.HandleError(c, err, "failed to create new order", http.StatusInternalServerError)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		utils.HandleError(c, err, "failed to marshal order", http.StatusInternalServerError)
		return
	}
	err = h.producer.PushOrderToQueue(ProducerTopic, ProducerPort, orderInBytes)
	if err != nil {
		utils.HandleError(c, err, "failed to push order to queue", http.StatusInternalServerError)
		return
	}

}
