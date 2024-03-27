package carry

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Carrie, error)
	GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error)
	GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error)
	Exists(ctx context.Context, carrieCode string) bool
	ExistsFK(ctx context.Context, localityCode string) bool
	Crear(ctx context.Context, c domain.Carrie) (int, error)
	//Update(ctx context.Context, c domain.Carrie) error
	//Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// retorna todos los  Carries
func (r *repository) GetAll(ctx context.Context) ([]domain.Carrie, error) {
	query := "SELECT id, cid, company_name, address, telephone, locality_id FROM carries;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var carriesG []domain.Carrie

	for rows.Next() {
		c := domain.Carrie{}
		_ = rows.Scan(&c.Id, &c.Cid, &c.Company_name, &c.Address, &c.Telephone, &c.Locality_id)
		carriesG = append(carriesG, c)
	}

	return carriesG, nil
}

// retorna la cantidad de Carries por cada Locality
func (r *repository) GetByLocality(ctx context.Context) ([]domain.CarrieLocality, error) {

	query := "SELECT  localities.id , localities.local_name, COUNT(carries.cid) as carries_count FROM carries " +
		"INNER JOIN localities ON carries.locality_id = localities.id GROUP BY carries.locality_id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var carriesG []domain.CarrieLocality

	for rows.Next() {
		c := domain.CarrieLocality{}
		_ = rows.Scan(&c.Locality_id, &c.Locality_name, &c.Cant_carries)
		carriesG = append(carriesG, c)
	}

	return carriesG, nil

}

// retorna la cantidad de Carries de una cada Locality determinada
func (r *repository) GetByLocalityID(ctx context.Context, id string) (domain.CarrieLocality, error) {
	query := "SELECT  localities.id , localities.local_name, COUNT(carries.cid) as carries_count FROM carries " +
		"INNER JOIN localities ON carries.locality_id = localities.id GROUP BY carries.locality_id HAVING carries.locality_id =?;"
	row := r.db.QueryRow(query, id)
	c := domain.CarrieLocality{}
	err := row.Scan(&c.Locality_id, &c.Locality_name, &c.Cant_carries)

	if err != nil {
		fmt.Println("error: ", err)
		return domain.CarrieLocality{}, err
	}

	return c, nil
}

func (r *repository) Exists(ctx context.Context, carrieCode string) bool {
	query := "SELECT cid FROM carries WHERE cid=?;"
	row := r.db.QueryRow(query, carrieCode)
	err := row.Scan(&carrieCode) //sino hay coincidencia retorna ErrNoRows
	return err == nil            //retorna true si existe id
}

func (r *repository) ExistsFK(ctx context.Context, localityCode string) bool {
	query := "SELECT id FROM localities WHERE id=?;"
	row := r.db.QueryRow(query, localityCode)
	err := row.Scan(&localityCode)
	return err == nil //retorna true si existe FK
}

func (r *repository) Crear(ctx context.Context, c domain.Carrie) (int, error) {
	query := "INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err //devuelve valor por defecto
	}

	res, err := stmt.Exec(&c.Cid, &c.Company_name, &c.Address, &c.Telephone, &c.Locality_id)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

/*func (r *repository) Update(ctx context.Context, c domain.Carrie) error {
	query := "UPDATE carries SET cid=?, company_name=?, address=?, telephone=?, locality_id=? WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&c.Cid, &c.Company_name, &c.Address, &c.Telephone, &c.Locality_id)
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
	query := "DELETE FROM carries WHERE id=?"
	stmt, err := r.db.Prepare(query)
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
}*/
