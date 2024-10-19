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

type MockOrderStore struct{}
type MockProductStore struct{}
type MockUserStore struct{}
type MockProducer struct{}

// Order store mocks
func (m *MockOrderStore) GetAllUsersOrders(id int) ([]types.Order, error) {
	return nil, nil
}
func (m *MockOrderStore) GetOrderByUniqueId(order types.Order) (types.Order, error) {
	return types.Order{}, nil
}
func (m *MockOrderStore) GetOrderByUserId(order types.Order) (*types.Order, error) {
	return nil, nil
}
func (m *MockOrderStore) CreateOrder(order types.Order) error {
	return nil
}
func (m *MockOrderStore) DeleteOrder(orderId string) error {
	return nil
}
func (m *MockOrderStore) GetUUIDFromOrder(orderId string) (string, error) {
	if orderId == "3832df47-0d0c-4aa5-bc2c-111111111111" {
		return "", fmt.Errorf("order is not ready yet")
	}
	return "", nil
}

// Product store mocks
func (m *MockProductStore) GetAllProducts() ([]types.Product, error) {
	return []types.Product{}, nil
}
func (m *MockProductStore) GetProductById(id int) (*types.Product, error) {
	if id > 10 && id < 20 {
		return &types.Product{
			Id:       10,
			Name:     "testName",
			Quantity: 1000,
			Price:    31,
		}, nil
	} else if id < 10 {
		return &types.Product{
			Id:       0,
			Name:     "",
			Quantity: 0,
			Price:    0,
		}, nil

	} else if id > 20 && id < 30 {
		return &types.Product{
			Id:       23,
			Name:     "",
			Quantity: 10,
			Price:    10,
		}, nil
	}
	return nil, fmt.Errorf("fafa")
}
func (m *MockProductStore) UpdateProductQuantity(id, orderQuantity, productQuantity int, action string) error {
	return nil
}

// User store mocks
func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}
func (m *MockUserStore) GetUserById(id float64) (*types.User, error) {
	return nil, nil
}
func (m *MockUserStore) CreateUser(user types.User) error {
	return nil
}

// kafka producer mock
func (m *MockProducer) PushOrderToQueue(topic, producerPort string, message []byte) error {
	return nil
}

func TestCreateOrder(t *testing.T) {
	orderStore := &MockOrderStore{}
	productStore := &MockProductStore{}
	userStore := &MockUserStore{}
	producer := &MockProducer{}
	Handler := NewHandler(orderStore, productStore, userStore, producer)
	t.Run("should pass if user payload is correct", func(t *testing.T) {
		var orders []*types.CreateOrder
		testUser := &types.User{
			Id: 1,
		}
		order1 := types.CreateOrder{
			Id:       11,
			Quantity: 1,
		}
		order2 := types.CreateOrder{
			Id:       12,
			Quantity: 2,
		}
		orders = append(orders, &order1)
		orders = append(orders, &order2)

		createOrderPayload := types.CreateOrderPayload{
			Orders: orders,
		}
		marshaledPayload, _ := json.Marshal(&createOrderPayload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {

			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/order", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.CreateOrder(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected : %d ,got : %d ", http.StatusOK, rr.Code)
		}
	})

	t.Run("Should fail if id is not found", func(t *testing.T) {
		var orders []*types.CreateOrder
		order1 := types.CreateOrder{
			Id:       1,
			Quantity: 10,
		}
		order2 := types.CreateOrder{
			Id:       2,
			Quantity: 12,
		}
		testUser := &types.User{
			Id: 1,
		}
		orders = append(orders, &order1)
		orders = append(orders, &order2)
		createOrderPayload := types.CreateOrderPayload{
			Orders: orders,
		}
		marshaledPayload, _ := json.Marshal(&createOrderPayload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/order", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.CreateOrder(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected : %d ,got : %d ", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("Should fail if quantity of order is bigger than quantity of product", func(t *testing.T) {
		var orders []*types.CreateOrder
		order1 := types.CreateOrder{
			Id:       23,
			Quantity: 1000,
		}
		order2 := types.CreateOrder{
			Id:       24,
			Quantity: 1000,
		}
		testUser := &types.User{
			Id: 1,
		}
		orders = append(orders, &order1)
		orders = append(orders, &order2)
		createOrderPayload := types.CreateOrderPayload{
			Orders: orders,
		}
		marshaledBody, _ := json.Marshal(&createOrderPayload)
		req, err := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(marshaledBody))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/order", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.CreateOrder(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected : %d , got : %d", http.StatusBadRequest, rr.Code)
		}

	})
	t.Run("Should pass Handler GetMyOrders", func(t *testing.T) {
		testUser := &types.User{
			Id: 1,
		}
		req, err := http.NewRequest(http.MethodGet, "/myorders", nil)
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.GET("myorders", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.GetMyOrders(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected: %d, got : %d", http.StatusOK, rr.Code)
		}

	})
	t.Run("Should pass if order exists", func(t *testing.T) {
		takeOrdersPayload := types.TakeOrderPayload{
			OrderId: []string{"3832df47-0d0c-4aa5-bc2c-36299bcd7f07"},
		}
		testUser := &types.User{
			Id: 1,
		}
		marshaledTakeORderPayload, _ := json.Marshal(&takeOrdersPayload)
		req, err := http.NewRequest(http.MethodPost, "/orders/take", bytes.NewBuffer(marshaledTakeORderPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/orders/take", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.GetMyOrders(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected : %d ,got : %d", http.StatusOK, rr.Code)
		}

	})
	t.Run("Should fail if order in not ready", func(t *testing.T) {
		takeOrdersPayload := types.TakeOrderPayload{
			OrderId: []string{"3832df47-0d0c-4aa5-bc2c-111111111111"},
		}
		testUser := &types.User{
			Id: 1,
		}
		marshaledTakeOrderPayload, _ := json.Marshal(&takeOrdersPayload)
		req, err := http.NewRequest(http.MethodPost, "/orders/take", bytes.NewBuffer(marshaledTakeOrderPayload))
		if err != nil {
			t.Log(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/orders/take", func(ctx *gin.Context) {
			ctx.Set("user", testUser)
			Handler.TakeOrders(ctx)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Log(takeOrdersPayload)
			t.Errorf("expected : %d ,got : %d", http.StatusBadRequest, rr.Code)
		}

	})

}
