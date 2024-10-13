package order

import (
	"net/http"

	"github.com/API/services/auth"
	"github.com/API/types"
	"github.com/API/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
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
	}

	err = h.store.CreateOrder(order)
	if err != nil {
		utils.HandleError(c, err, "failed to create new order", http.StatusInternalServerError)
		return
	}
}
