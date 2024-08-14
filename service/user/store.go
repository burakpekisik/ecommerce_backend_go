package user

import (
	"database/sql"
	"fmt"

	"github.com/burakpekisik/ecommerce_backend_go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders WHERE userId = ?", userID)

	if err != nil {
		return nil, err
	}

	orders := []types.Order{}
	for rows.Next() {
		o := new(types.Order)
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Total,
			&o.Status,
			&o.Address,
			&o.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, *o)
	}

	return orders, nil
}
