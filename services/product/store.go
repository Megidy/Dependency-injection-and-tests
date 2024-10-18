package product

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

func (s *Store) GetAllProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.Id,
		&product.Name,
		&product.Quantity,
		&product.Price,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) GetProductById(id int) (*types.Product, error) {
	var product types.Product
	row, err := s.db.Query("select * from products where id =?", id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.Price)
		if err != nil {
			return nil, err
		}
	}

	return &product, nil
}

func (s *Store) UpdateProductQuantity(id, orderQuantity, productQuantity int, action string) error {
	if action == "dec" {
		_, err := s.db.Exec("update products set quantity =? where id =?", productQuantity-orderQuantity, id)
		if err != nil {
			return err
		}
	} else if action == "inc" {
		_, err := s.db.Exec("update products set quantity =? where id =?", productQuantity+orderQuantity, id)
		if err != nil {
			return err
		}
	}
	return nil
}
