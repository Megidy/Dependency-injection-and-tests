package order

import (
	"database/sql"

	"github.com/API/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateOrder(order types.Order) error {
	_, err := s.db.Exec("insert into orders values(?,?,?)", order.UserID, order.Product.Id, order.Status)

	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetOrder(order types.Order) (*types.Order, error) {
	row, err := s.db.Query("select * from orders where user_id = ?", order.UserID)
	if err != nil {
		return nil, err

	}
	var o types.Order
	for row.Next() {
		err := row.Scan(&o.UserID, &o.Product.Id, &o.Product.Name, &o.Product.Quantity, &o.Product.Price)
		if err != nil {
			return nil, err
		}
	}
	return &o, nil
}
