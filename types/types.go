package types

import "github.com/google/uuid"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id float64) (*User, error)
	CreateUser(user User) error
}
type Producer interface {
	PushOrderToQueue(topic, producerPort string, message []byte) error
}
type Consumer interface {
	ReceiveOrders()
}

// TODO : ADD GET ORDER
type OrderStore interface {
	GetAllUsersOrders(id int) ([]Order, error)
	GetOrderByUniqueId(order Order) (Order, error)
	GetOrderByUserId(order Order) (*Order, error)
	CreateOrder(order Order) error
}

type ProductStore interface {
	GetAllProducts() ([]Product, error)
	GetProductById(id int) (*Product, error)
	UpdateProductQuantity(id, orderQuantity, productQuantity int, action string) error
}

type DepotStore interface {
	UpdateOrderStatus(order Order) error
}

type Product struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignInPayload struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LogInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO : ADD STATUS
type Order struct {
	Id       uuid.UUID
	UserID   int
	Product  Product
	Quantity int
	Status   string
}
type CreateOrder struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
type CreateOrderPayload struct {
	Orders []*CreateOrder
}
