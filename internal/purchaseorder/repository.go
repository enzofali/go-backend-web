package purchaseorder

import (
	"context"
	"database/sql"

	//"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, purchOrd domain.Purchase_Orders) (int, error)
	Exists(ctx context.Context, id int) bool
	ExistsBuyer(ctx context.Context, id int) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, purchOrd domain.Purchase_Orders) (int, error) {
	query := "INSERT INTO purchase_orders(order_number, order_date, tracking_code, buyer_id, product_record_id, order_status_id) VALUES (?,?,?,?,?,?);"

	statement, err := r.db.Prepare(query)

	if err != nil {
		return 0, ErrDatabase
	}

	res, err := statement.Exec(&purchOrd.Order_number,
		&purchOrd.Order_date,
		&purchOrd.Tracking_code,
		&purchOrd.Buyer_id,
		&purchOrd.Product_record_id,
		&purchOrd.Order_Status_id)

	if err != nil {
		errMysql, err2 := err.(*mysql.MySQLError)
		if !err2 {
			return 0, ErrDatabase
		}
		switch errMysql.Number {
		case 1062:
			return 0, ErrExists
		default:
			return 0, ErrDatabase
		}
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, ErrDatabase
	}
	purchOrd.ID = int(id)
	return int(id), nil
}

func (r *repository) Exists(ctx context.Context, id int) bool {
	query := "SELECT id FROM purchase_orders WHERE id=?"
	/*statement, err := r.db.Prepare(query)

	if err!= nil{
		return false
	}*/
	row := r.db.QueryRow(query, id)
	exist := row.Scan(&id)
	return exist == nil
}

func (r *repository) ExistsBuyer(ctx context.Context, id int) bool {
	query := "SELECT id FROM buyers WHERE id=?"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}
