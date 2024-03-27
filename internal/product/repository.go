package product

import (
	"context"
	"database/sql"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

// Repository encapsulates the storage of a Product.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Exists(ctx context.Context, productCode string) bool
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, p domain.Product) error
	Delete(ctx context.Context, id int) error
	ValidateProductID(ctx context.Context, pid int) bool
	GetOneReport(ctx context.Context, id int) (int, string, error)
	GetAllReports(ctx context.Context) ([]domain.Report, error)
	StoreType(ctx context.Context, name string) (int, error)
}

type repository struct {
	db *sql.DB
}

// Queries
var (
	GET_ALL = `
		SELECT
			id,description,expiration_rate,freezing_rate,height,lenght,netweight,product_code,recommended_freezing_temperature,width,id_product_type,id_seller
		FROM
			products;
	`
	GET_ONE = `
		SELECT
			id,description,expiration_rate,freezing_rate,height,lenght,netweight,product_code,recommended_freezing_temperature,width,id_product_type,id_seller 
		FROM
			products
		WHERE
			id=?;
	`
	EXISTS = `SELECT product_code FROM products WHERE product_code=?;`
	SAVE   = `
		INSERT INTO 
			products(description,expiration_rate,freezing_rate,height,lenght,netweight,product_code,recommended_freezing_temperature,width,id_product_type,id_seller)
		VALUES
			(?,?,?,?,?,?,?,?,?,?,?);
	`
	UPDATE = `
		UPDATE
			products
		SET 
			description=?, expiration_rate=?, freezing_rate=?, height=?, lenght=?, netweight=?, product_code=?, recommended_freezing_temperature=?, width=?, id_product_type=?, id_seller=? 
		WHERE
			id=?;
	`
	DELETE         = `DELETE FROM products WHERE id=?;`
	VALIDATE       = `SELECT COUNT(id) FROM products WHERE id = ?;`
	GET_ONE_REPORT = `
		SELECT COUNT(pr.id), p.description FROM product_records pr
		INNER JOIN products p ON pr.product_id = p.id
		WHERE pr.product_id = ?
		GROUP BY p.id;
	`
	GET_ALL_REPORTS = `
		SELECT COUNT(pr.id), p.description, p.id FROM product_records pr
		LEFT JOIN products p ON pr.product_id = p.id
		GROUP BY p.id;
	`
	STORE_TYPE = `INSERT INTO product_types(name) VALUES (?);`
)

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Product, error) {
	stmt, err := r.db.Prepare(GET_ALL)
	if err != nil {
		return []domain.Product{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return []domain.Product{}, err
	}

	var products []domain.Product

	for rows.Next() {
		p := domain.Product{}
		_ = rows.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.Netweight, &p.ProductCode, &p.RecomFreezTemp, &p.Width, &p.ProductTypeID, &p.SellerID)
		products = append(products, p)
	}

	return products, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Product, error) {
	stmt, err := r.db.Prepare(GET_ONE)
	if err != nil {
		return domain.Product{}, err
	}

	row := stmt.QueryRow(id)
	p := domain.Product{}

	err = row.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.Netweight, &p.ProductCode, &p.RecomFreezTemp, &p.Width, &p.ProductTypeID, &p.SellerID)
	if err != nil {
		return domain.Product{}, err
	}

	return p, nil
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	row := r.db.QueryRow(EXISTS, productCode)
	err := row.Scan(&productCode)
	return err == nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (int, error) {
	stmt, err := r.db.Prepare(SAVE)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.Netweight, p.ProductCode, p.RecomFreezTemp, p.Width, p.ProductTypeID, p.SellerID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, p domain.Product) error {
	stmt, err := r.db.Prepare(UPDATE)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.Netweight, p.ProductCode, p.RecomFreezTemp, p.Width, p.ProductTypeID, p.SellerID, p.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.Prepare(DELETE)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}

// Product Type
// creates a new product type and returns its id
func (r *repository) StoreType(ctx context.Context, name string) (int, error) {
	stmt, err := r.db.Prepare(STORE_TYPE)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Product Record Reports

// Checks whether a given id exists and is unique in the products table
func (r *repository) ValidateProductID(ctx context.Context, pid int) bool {
	var count int
	stmt, err := r.db.Prepare(VALIDATE)
	if err != nil {
		return false
	}
	row := stmt.QueryRow(pid)
	err = row.Scan(&count)
	// id is valid if there are no errors and the count is exactly 1
	return err == nil && count == 1
}

// Returns the number of records in the product_records table associated with the given product_id,
// as well as the product's description.
func (r *repository) GetOneReport(ctx context.Context, id int) (int, string, error) {
	var count int
	var description string
	stmt, err := r.db.Prepare(GET_ONE_REPORT)
	if err != nil {
		return 0, "", err
	}
	row := stmt.QueryRow(id)
	err = row.Scan(&count, &description)
	if err != nil {
		return 0, "", err
	}

	return count, description, nil
}

// Returns the number of records in the product_records table for each product.
func (r *repository) GetAllReports(ctx context.Context) ([]domain.Report, error) {
	var reports []domain.Report
	stmt, err := r.db.Prepare(GET_ALL_REPORTS)
	if err != nil {
		return []domain.Report{}, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return []domain.Report{}, err
	}

	for rows.Next() {
		var report domain.Report
		if err := rows.Scan(&report.Count, &report.Description, &report.ProductID); err != nil {
			return []domain.Report{}, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}
