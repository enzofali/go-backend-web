package product_records

import (
	"context"
	"database/sql"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

type Repository interface {
	ValidateProductID(ctx context.Context, pid int) bool
	Store(ctx context.Context, pr domain.ProductRecord) (int, error)
}

type repository struct {
	db *sql.DB
}

// queries
var (
	VALIDATE = `SELECT COUNT(id) FROM products WHERE id = ?`
	STORE    = `
		INSERT INTO product_records(last_update_date, purchase_price, sale_price, product_id)
		VALUES (?, ?, ?, ?)
	`
)

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ValidateProductID(ctx context.Context, pid int) bool {
	var count int
	stmt, err := r.db.Prepare(VALIDATE)
	if err != nil {
		return false
	}
	row := stmt.QueryRow(pid)
	err = row.Scan(&count)
	return err == nil && count == 1
}

func (r *repository) Store(ctx context.Context, pr domain.ProductRecord) (int, error) {
	stmt, err := r.db.Prepare(STORE)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(pr.LastUpdateDate, pr.PurchasePrice, pr.SalePrice, pr.ProductID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
