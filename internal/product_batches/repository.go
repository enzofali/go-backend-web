package product_batches

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

var (
	ErrExistsBatchNumber = errors.New("error: batch number already exists")
	ErrProductNotFound   = errors.New("error: product id does not exists")
	ErrSectionNotFound   = errors.New("error: section id does not exists")
	ErrInternal          = errors.New("error: internal error")
)

var (
	createQuery = "INSERT INTO products_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
)

type Repository interface {
	Create(ctx context.Context, p domain.ProductBatches) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, p domain.ProductBatches) (int, error) {
	stmt, err := r.db.Prepare(createQuery)
	if err != nil {
		return 0, ErrInternal
	}
	defer stmt.Close()

	res, err := stmt.Exec(&p.BatchNumber, &p.CurrentQuantity, &p.CurrentTemperature, &p.DueDate, &p.InitialQuantity, &p.ManufacturingDate, &p.ManufacturingHour, &p.MinumumTemperature, &p.ProductID, &p.SectionID)

	if err != nil {
		// log.Println(err.(*mysql.MySQLError).Number, err.(*mysql.MySQLError).Message)
		errMysql := err.(*mysql.MySQLError)
		switch errMysql.Number {
		case 1452:
			if strings.Contains(errMysql.Error(), "`products`") {
				return 0, ErrProductNotFound
			} else if strings.Contains(errMysql.Error(), "`sections`") {
				return 0, ErrSectionNotFound
			}
		case 1062:
			return 0, ErrExistsBatchNumber
		default:
			return 0, ErrInternal
		}
	}

	rows, err := res.RowsAffected()
	if err != nil || rows != 1 {
		return 0, ErrInternal
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, ErrInternal
	}

	return int(id), nil
}
