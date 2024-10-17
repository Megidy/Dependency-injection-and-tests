package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/API/types"
	"github.com/gin-gonic/gin"
)

type MockProducer struct{}
type MockOrderStore struct{}
type MockUserStore struct{}
type MockProductStore struct{}

func (m *MockOrderStore) CreateOrder(order types.Order) error {
	if order.Product.Id == 60 {
		return fmt.Errorf("failed to create order")
	}
	return nil
}

func (m *MockOrderStore) GetOrder(order types.Order) (*types.Order, error) {
	return nil, nil

}
func (m *MockUserStore) GetUserById(id float64) (*types.User, error) {
	return nil, nil
}
func (m *MockUserStore) CreateUser(user types.User) error {
	return nil
}
func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *MockProductStore) GetAllProducts() ([]*types.Product, error) {
	return nil, nil
}
func (m *MockProductStore) GetProductById(id int) (*types.Product, error) {
	if id == 1 {
		return &types.Product{
			Id:       1,
			Quantity: 1,
		}, nil
	} else if id == 0 {
		return &types.Product{
			Id:       0,
			Quantity: 1,
		}, nil
	} else if id == 2 {
		return &types.Product{
			Id:       2,
			Quantity: 1,
		}, nil

	}

	return nil, nil
}

func (m *MockProducer) PushOrderToQueue(topic, producerPort string, message []byte) error {
	return nil
}
func TestCreateOrder(t *testing.T) {
	newProducer := &MockProducer{}
	orderStore := &MockOrderStore{}
	productStore := &MockProductStore{}
	userStore := &MockUserStore{}
	orderHandler := NewHandler(orderStore, productStore, userStore, newProducer)

	t.Run("Should pass if user payload is correct", func(t *testing.T) {

		payload := types.CreateOrderPayload{
			Id:       1,
			Quantity: 1,
		}

		testUser := &types.User{
			Id: 1,
		}

		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/order", func(c *gin.Context) {

			c.Set("user", testUser)
			orderHandler.CreateOrder(c)
		})

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
		}
	})
	t.Run("Should fail if id of product is not avaible", func(t *testing.T) {
		payload := types.CreateOrderPayload{
			Id:       0,
			Quantity: 1,
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/order", orderHandler.CreateOrder)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected %d,got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("Should fail if quantity is higher than avaible", func(t *testing.T) {
		payload := types.CreateOrderPayload{
			Id:       1,
			Quantity: 300,
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/order", orderHandler.CreateOrder)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected %d,got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("Should fail if didnt create new order", func(t *testing.T) {
		payload := types.CreateOrderPayload{
			Id:       60,
			Quantity: 1,
		}
		testUser := &types.User{
			Id: 1,
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/order", func(c *gin.Context) {
			c.Set("user", testUser)
			orderHandler.CreateOrder(c)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
