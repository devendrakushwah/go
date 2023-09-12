package testcontainers

import "database/sql"

type OrderService struct {
	db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) AddOrder(productID int, quantity int) error {
	_, err := s.db.Exec("INSERT INTO orders (product_id, quantity) VALUES (?, ?)", productID, quantity)
	return err
}

func (s *OrderService) GetOrder(id int) (int, int, error) {
	var productID, quantity int
	err := s.db.QueryRow("SELECT product_id, quantity FROM orders WHERE id=?", id).Scan(&productID, &quantity)
	if err != nil {
		return 0, 0, err
	}
	return productID, quantity, nil
}
