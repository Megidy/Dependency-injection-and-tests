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
	_, err := s.db.Exec("insert into orders values(?,?)", order.UserID, order.Product.Id)

	if err != nil {
		return err
	}
	return nil
}
