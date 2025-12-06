package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrOutofStock = errors.New("product out of stock")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// PurchaseProduct handles the ACID transaction
func (s *Store) PurchaseProduct(ctx context.Context, userID, productID, quantity int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// 2. LOCK THE ROW (Pessimistic Locking)
	var currentStock int
	queryCheck := `SELECT stock FROM products WHERE id = $1 FOR UPDATE`

	err = tx.QueryRowContext(ctx, queryCheck, productID).Scan(&currentStock)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	// 3. Check Stock Logic
	if currentStock < quantity {
		return ErrOutofStock
	}

	// 4. Deduct Stock
	queryUpdate := `UPDATE products SET stock = stock - $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, queryUpdate, quantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// 5. Create Order Record
	queryInsert := `INSERT INTO orders (user_id, product_id, quantity) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, queryInsert, userID, productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// 6. Commit the Transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("✅ Order placed: User %d bought Product %d", userID, productID)
	return nil
}

// Helper to reset stock (useful for testing repeated runs)
func (s *Store) ResetStock(productID, stock int) error {
	_, err := s.db.Exec("UPDATE products SET stock = $1 WHERE id = $2", stock, productID)
	return err
}
