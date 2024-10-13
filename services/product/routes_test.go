package product

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/API/types"
	"github.com/gin-gonic/gin"
)

type MockProductStore struct {
}

func (m *MockProductStore) GetAllProducts() ([]*types.Product, error) {
	return nil, nil
}
func (m *MockProductStore) GetProductById(id int) (*types.Product, error) {
	return nil, nil
}

func TestGetAllProducts(t *testing.T) {

	productStore := &MockProductStore{}
	productHandler := NewHandler(productStore)

	t.Run("Should handle get products ", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/products", productHandler.GetAllProducts)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
		}
	})

}
