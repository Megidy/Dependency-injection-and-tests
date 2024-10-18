package order

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/API/services/auth"
	"github.com/API/types"
	"github.com/API/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	router.GET("/myorders", auth.NewHandler(h.userStore).WithJWT, h.GetMyOrders)
}
func (h *Handler) CreateOrder(c *gin.Context) {
	user, _ := c.Get("user")
	var payload types.CreateOrderPayload
	var order types.Order
	var products []types.Product
	var quant int
	var response []types.Order
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		utils.HandleError(c, err, "failed to read body", http.StatusBadRequest)
		return
	}

	for _, order := range payload.Orders {
		product, err := h.productStore.GetProductById(order.Id)
		if err != nil {
			utils.HandleError(c, err, "failed to get sinle product", http.StatusBadRequest)
			return
		}

		if product.Id == 0 {
			utils.HandleError(c, err, "no product with this id", http.StatusBadRequest)
			return
		}
		if order.Quantity > product.Quantity {
			utils.HandleError(c, err, "not avaible in this amount", http.StatusBadRequest)
			return
		}
		err = h.productStore.UpdateProductQuantity(product.Id, order.Quantity, product.Quantity, "dec")
		if err != nil {
			utils.HandleError(c, err, "error when updating product quantity", http.StatusInternalServerError)
			return
		}
		quant = order.Quantity
		log.Println(product)
		products = append(products, *product)
	}
	for _, product := range products {
		order = types.Order{
			Id:       uuid.New(),
			UserID:   user.(*types.User).Id,
			Product:  product,
			Quantity: quant,
			Status:   "Pending",
		}
		response = append(response, order)
		log.Println(order)
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
	c.JSON(http.StatusOK, gin.H{
		"success": response,
	})
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	user, _ := c.Get("user")
	log.Println("got user")
	var orders []types.Order
	orders, err := h.store.GetAllUsersOrders(user.(*types.User).Id)
	log.Println("created orders")
	if err != nil {
		utils.HandleError(c, err, err.Error(), http.StatusBadRequest)
		return
	}
	for _, order := range orders {
		product, err := h.productStore.GetProductById(order.Product.Id)
		log.Println("order :", order)
		log.Println("prdocut :", product)
		if err != nil {
			utils.HandleError(c, err, "failed to get products", http.StatusBadRequest)
			return
		}

		order = types.Order{
			Product: *product,
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"your orders": orders,
	})
}
