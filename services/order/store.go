package order

import (
	"database/sql"
	"fmt"
	"log"

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
	_, err := s.db.Exec("insert into orders values(?,?,?,?,?)", order.Id, order.UserID, order.Product.Id, order.Quantity, order.Status)

	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetOrderByUserId(order types.Order) (*types.Order, error) {
	row, err := s.db.Query("select * from orders where user_id = ?", order.UserID)
	if err != nil {
		return nil, err

	}
	var o types.Order
	for row.Next() {
		err := row.Scan(&o.Id, &o.UserID, &o.Product.Id, &o.Product.Name, &o.Product.Quantity, &o.Product.Price)
		if err != nil {
			return nil, err
		}
	}
	return &o, nil
}
func (s *Store) GetOrderByUniqueId(order types.Order) (types.Order, error) {
	row, err := s.db.Query("select * from orders where id = ?", order.Id)
	if err != nil {
		return types.Order{}, err

	}
	var o types.Order
	for row.Next() {
		err := row.Scan(&o.Id, &o.UserID, &o.Product.Id, &o.Product.Name, &o.Product.Quantity, &o.Product.Price, &o.Status, &o.Quantity)
		if err != nil {
			return types.Order{}, err
		}
	}
	return o, nil
}
func (s *Store) GetAllUsersOrders(id int) ([]types.Order, error) {
	var orders []types.Order
	rows, err := s.db.Query("select id,user_id,quantity,status from orders where user_id=?", id)
	if err != nil {
		return nil, err
	}
	log.Println("queried")
	for rows.Next() {
		var order types.Order
		err = rows.Scan(&order.Id, &order.UserID, &order.Quantity, &order.Status)
		log.Println("Scanned")
		if err != nil {
			return nil, err

		}
		log.Println("appended")
		orders = append(orders, order)
	}
	log.Println("returned ")
	return orders, nil
}

func (s *Store) DeleteOrder(orderId string) error {
	log.Println(orderId)

	_, err := s.db.Exec("delete from orders where id=?", orderId)
	log.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUUIDFromOrder(orderId string) (string, error) {
	row, err := s.db.Query("select id,status from orders where id =?", orderId)
	if err != nil {
		return "", err
	}
	var id, status string
	for row.Next() {
		err = row.Scan(&id, &status)
		if status != "Ready to pickup" {
			return "", fmt.Errorf("order is not ready yet")
		}
		if err != nil {
			return "", err
		}
	}
	return id, nil

}
