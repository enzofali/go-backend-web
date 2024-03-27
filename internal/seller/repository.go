package seller

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

var (
	ErrIntern          = errors.New("an internal error")
	ErrDuplicated      = errors.New("duplicated locality")
	ErrInvalidLocality = errors.New("invalid id locality")
	ErrNotFound        = errors.New("seller not found")
	QueryGetAll        = "SELECT id,cid,company_name,address,telephone,locality_id FROM sellers"
	QueryGetById       = "SELECT id,cid,company_name,address,telephone,locality_id FROM sellers WHERE id=?;"
	QueryExistsCid     = "SELECT cid FROM sellers WHERE cid=?;"
	QueryInsert        = "INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	QueryUpdate        = "UPDATE sellers SET cid=?, company_name=?, address=?, telephone=?, locality_id=? WHERE id=?"
	QueryDelete        = "DELETE FROM sellers WHERE id=?"
)

// Repository encapsulates the storage of a Seller.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Seller) (int, error)
	Update(ctx context.Context, s domain.Seller) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Seller, error) {
	rows, err := r.db.Query(QueryGetAll)
	if err != nil {
		return nil, ErrIntern
	}

	var sellers []domain.Seller

	for rows.Next() {
		s := domain.Seller{}
		_ = rows.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.Locality_id)
		sellers = append(sellers, s)
	}

	return sellers, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Seller, error) {
	row := r.db.QueryRow(QueryGetById, id)
	s := domain.Seller{}
	err := row.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.Locality_id)
	if err != nil {

		switch err {
		case sql.ErrNoRows:
			err = ErrNotFound
		default:
			err = ErrIntern
		}
		return domain.Seller{}, err
	}

	return s, nil
}

func (r *repository) Exists(ctx context.Context, cid int) bool {
	row := r.db.QueryRow(QueryExistsCid, cid)
	err := row.Scan(&cid)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Seller) (int, error) {
	stmt, err := r.db.Prepare(QueryInsert)
	if err != nil {
		return 0, ErrIntern
	}

	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id)
	if err != nil {
		driverErr, ok := err.(*mysql.MySQLError)
		if !ok {
			err = ErrIntern
			return 0, err
		}

		switch driverErr.Number {
		case 1452:
			err = ErrInvalidLocality
		default:
			err = ErrIntern
		}

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, ErrIntern
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, s domain.Seller) error {
	stmt, err := r.db.Prepare(QueryUpdate)
	if err != nil {
		return ErrIntern
	}

	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.Locality_id, s.ID)
	if err != nil {
		return ErrIntern
	}

	_, err = res.RowsAffected()
	if err != nil {
		return ErrIntern
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(QueryDelete)
	if err != nil {
		return ErrIntern
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return ErrIntern
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return ErrIntern
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}
