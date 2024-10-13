package types

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id float64) (*User, error)
	CreateUser(user User) error
}

type ProductStore interface {
	GetAllProducts() ([]*Product, error)
	GetProductById(id int) (*Product, error)
}

type Product struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type Order struct {
	UserID  int
	Product *Product
}
type OrderStore interface {
	CreateOrder(Order) error
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

type CreateOrderPayload struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
