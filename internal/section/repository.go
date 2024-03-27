package section

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Errors
var (
	ErrInternal            = errors.New("error: internal error")
	ErrExistsSectionNumber = errors.New("error: section number already exists")
	ErrWareHouseNotFound   = errors.New("error: warehouse id does not exists")
	ErrProductTypeNotFound = errors.New("error: product type id does not exists")
	ErrSectionNotFound     = errors.New("error: section id does not exists")
)

var (
	GetAllQuery = "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections;"
	GetByID     = "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type FROM sections WHERE id=?;"
	// Products quantity of each Section
	GetReportQuery = "SELECT s.id, s.section_number, COALESCE(sum(pb.current_quantity),0) FROM sections as s " +
		"LEFT JOIN products_batches as pb ON s.id = pb.section_id " +
		"GROUP BY s.id, s.section_number;"
	// Products quantity for a certain section
	GetReportQueryByID = "SELECT s.id, s.section_number, COALESCE(sum(pb.current_quantity),0) FROM sections as s " +
		"LEFT JOIN products_batches as pb ON s.id = pb.section_id " +
		"WHERE s.id = ? " +
		"GROUP BY s.id, s.section_number;"
	CreateQuery = "INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
	UpdateQuery = "UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, id_product_type=? WHERE id=?;"
	DeleteQuery = "DELETE FROM sections WHERE id=?;"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	GetByID(ctx context.Context, id int) (domain.Section, error)
	GetAllReportProducts(ctx context.Context) ([]domain.SectionReportProducts, error)
	GetReportProductsByID(ctx context.Context, id int) ([]domain.SectionReportProducts, error)
	Create(ctx context.Context, s domain.Section) (int, error)
	Update(ctx context.Context, s domain.Section) error
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

// ------------------------------- READ ---------------------------------

func (r *repository) GetAll(ctx context.Context) ([]domain.Section, error) {
	rows, err := r.db.Query(GetAllQuery)
	if err != nil {
		return nil, ErrInternal
	}

	var sections []domain.Section

	for rows.Next() {
		s := domain.Section{}
		err = rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		if err != nil {
			return []domain.Section{}, ErrInternal
		}
		sections = append(sections, s)
	}

	return sections, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (domain.Section, error) {
	row := r.db.QueryRow(GetByID, id)
	s := domain.Section{}
	err := row.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.Section{}, ErrSectionNotFound
		default:
			return domain.Section{}, ErrInternal
		}
	}

	return s, nil
}

func (r *repository) GetAllReportProducts(ctx context.Context) ([]domain.SectionReportProducts, error) {
	var reports []domain.SectionReportProducts

	rows, err := r.db.Query(GetReportQuery)
	if err != nil {
		return []domain.SectionReportProducts{}, ErrInternal
	}

	for rows.Next() {
		report := domain.SectionReportProducts{}
		err := rows.Scan(&report.ID, &report.SectionNumber, &report.ProductCount)
		if err != nil {
			return []domain.SectionReportProducts{}, ErrInternal
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *repository) GetReportProductsByID(ctx context.Context, id int) ([]domain.SectionReportProducts, error) {
	var reports []domain.SectionReportProducts

	row := r.db.QueryRow(GetReportQueryByID, id)
	report := domain.SectionReportProducts{}
	err := row.Scan(&report.ID, &report.SectionNumber, &report.ProductCount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []domain.SectionReportProducts{}, ErrSectionNotFound
		default:
			return []domain.SectionReportProducts{}, ErrInternal
		}
	}

	reports = append(reports, report)
	return reports, nil
}

// -------------------------------- WRITE --------------------------------

func (r *repository) Create(ctx context.Context, s domain.Section) (int, error) {
	stmt, err := r.db.Prepare(CreateQuery)
	if err != nil {
		return 0, ErrInternal
	}
	defer stmt.Close()

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		// log.Println(err.(*mysql.MySQLError).Number)
		// log.Println(err.(*mysql.MySQLError).Message)
		errMysql := err.(*mysql.MySQLError)
		switch errMysql.Number {
		case 1452:
			if strings.Contains(errMysql.Error(), "`warehouses`") {
				return 0, ErrWareHouseNotFound
			} else if strings.Contains(errMysql.Error(), "`product_types`") {
				return 0, ErrProductTypeNotFound
			}
		case 1062:
			return 0, ErrExistsSectionNumber
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

func (r *repository) Update(ctx context.Context, s domain.Section) error {
	stmt, err := r.db.Prepare(UpdateQuery)
	if err != nil {
		return ErrInternal
	}
	defer stmt.Close()

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID, &s.ID)
	if err != nil {
		// log.Println(err.(*mysql.MySQLError).Number, err.(*mysql.MySQLError).Message)
		errMysql := err.(*mysql.MySQLError)
		switch errMysql.Number {
		case 1452:
			if strings.Contains(errMysql.Error(), "`warehouses`") {
				return ErrWareHouseNotFound
			} else if strings.Contains(errMysql.Error(), "`product_types`") {
				return ErrProductTypeNotFound
			}
		case 1062:
			return ErrExistsSectionNumber
		default:
			return ErrInternal
		}
	}

	_, err = res.RowsAffected()
	if err != nil {
		return ErrInternal
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {

	stmt, err := r.db.Prepare(DeleteQuery)
	if err != nil {
		return ErrInternal
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return ErrInternal
	}

	affect, err := res.RowsAffected()
	if err != nil || affect < 1 {
		return ErrInternal
	}

	return nil
}
