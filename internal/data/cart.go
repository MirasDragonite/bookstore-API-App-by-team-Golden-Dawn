package data

import (
	"context"
	"database/sql"
	"time"
)

type Cart struct {
	ID int64 `json:"id"`

	BookID int64 `json:"bookID"`
}

type CartModel struct {
	DB *sql.DB
}

func (m CartModel) Ins(cart *Cart) error {

	query := `
    INSERT INTO cart (bookID)
    VALUES ($1)
    RETURNING id`
	args := []any{cart.BookID}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&cart.ID)
}

func (m CartModel) RemoveFromCart(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
DELETE FROM cart
WHERE id = $1`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m CartModel) Order() error {

	query := `
DELETE * FROM cart
`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
