package order

import (
	"database/sql"
	"fmt"

	"github.com/burakpekisik/ecommerce_backend_go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	fmt.Println(id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)

	return err
}

func (s *Store) CancelOrder(orderID int) error {
	var status string

	// Siparişin durumunu veritabanından al
	err := s.db.QueryRow("SELECT status FROM orders WHERE id = ?", orderID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		return err
	}

	// Eğer sipariş zaten iptal edilmişse hata döndür
	if status == "cancelled" {
		return fmt.Errorf("order is already cancelled")
	} else if status == "shipped" {
		return fmt.Errorf("order is already shipped")
	}

	// Siparişi iptal et
	_, err = s.db.Exec("UPDATE orders SET status = 'cancelled' WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	// Her bir ürünün miktarını artır
	rows, err := s.db.Query("SELECT productId, quantity FROM order_items WHERE orderId = ?", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productId int
		var quantity int

		if err := rows.Scan(&productId, &quantity); err != nil {
			return err
		}

		_, err = s.db.Exec("UPDATE products SET quantity = quantity + ? WHERE id = ?", quantity, productId)
		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// Siparişe ait ürünleri sil
	_, err = s.db.Exec("DELETE FROM order_items WHERE orderId = ?", orderID)
	if err != nil {
		return err
	}

	return nil
}
