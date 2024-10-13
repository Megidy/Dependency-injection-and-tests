package product

import (
	"net/http"

	"github.com/API/types"
	"github.com/API/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/products", h.GetAllProducts)
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		utils.HandleError(c, err, "failed to get products", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"All products": products,
	})

}
