package inboundorder

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

var (
	ErrEmployeeNotFound     = errors.New("Employee not found")
	ErrProductBatchNotFound = errors.New("Product Batch not found")
	ErrWarehouseNotFound    = errors.New("Warehouse not found")
	ErrOrderNumberExtists   = errors.New("Order number exists")
)

type Repository interface {
	Save(ctx context.Context, i domain.InboundOrder) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, i domain.InboundOrder) (int, error) {
	query := "INSERT INTO inbound_orders(order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?,?,?,?,?)"
	smt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer smt.Close()

	res, err := smt.Exec(&i.OrderDate, &i.OrderNumber, &i.EmployeeID, &i.ProductBatchID, &i.WarehouseID)

	if err != nil {
		log.Println(err.(*mysql.MySQLError).Number)
		log.Println(err.(*mysql.MySQLError).Message)
		errMysql := err.(*mysql.MySQLError)
		switch errMysql.Number {
		case 1452:
			if strings.Contains(errMysql.Error(), "employees") {
				return 0, ErrEmployeeNotFound
			}
			if strings.Contains(errMysql.Error(), "products_batches") {
				return 0, ErrProductBatchNotFound
			}
			if strings.Contains(errMysql.Error(), "warehouses") {
				return 0, ErrWarehouseNotFound
			}
		case 1062:
			return 0, ErrOrderNumberExtists
		default:
			return 0, err
		}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
