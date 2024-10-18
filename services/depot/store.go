package depot

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

func (s *Store) UpdateOrderStatus(order types.Order) error {
	_, err := s.db.Exec("update orders set status =?  where id =?", order.Status, order.Id)
	if err != nil {
		return err
	}
	return nil
}
