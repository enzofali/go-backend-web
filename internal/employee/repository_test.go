package employee

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Repository_GetAll(t *testing.T) {
	ctx := context.Background()

	query := "SELECT * FROM employees"

	data := []domain.Employee{
		{
			ID:           1,
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
			CardNumberID: "12",
		},
		{
			ID:           2,
			FirstName:    "Carlos",
			LastName:     "Perez",
			WarehouseID:  1,
			CardNumberID: "121",
		},
	}
	t.Run("GetAll OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		for _, d := range data {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName, d.WarehouseID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := NewRepository(db)
		employees, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, data, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAll Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)
		employees, err := repo.GetAll(ctx)
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAll Error Scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})
		for _, d := range data {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := NewRepository(db)
		employees, err := repo.GetAll(ctx)
		assert.Error(t, err)
		assert.Empty(t, employees, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_Get(t *testing.T) {
	ctx := context.Background()

	query := "SELECT * FROM employees WHERE id=?;"

	data := domain.Employee{
		ID:           1,
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
		CardNumberID: "12",
	}

	t.Run("Get OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		rows.AddRow(data.ID, data.CardNumberID, data.FirstName, data.LastName, data.WarehouseID)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

		repo := NewRepository(db)

		employee, err := repo.Get(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, data, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)

		employee, err := repo.Get(ctx, 1)
		assert.Error(t, err)
		assert.Empty(t, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_Exists(t *testing.T) {
	ctx := context.Background()

	query := "SELECT card_number_id FROM employees WHERE card_number_id=?;"

	data := domain.Employee{
		ID:           1,
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
		CardNumberID: "12",
	}

	t.Run("Exists True", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"card_number_id"})
		rows.AddRow(data.CardNumberID)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("12").WillReturnRows(rows)

		repo := NewRepository(db)

		exists := repo.Exists(ctx, "12")
		assert.Equal(t, true, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exists False", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"card_number_id"})

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("12").WillReturnRows(rows)

		repo := NewRepository(db)

		exists := repo.Exists(ctx, "12")
		assert.Equal(t, false, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_Save(t *testing.T) {
	ctx := context.Background()

	query := "INSERT INTO employees(card_number_id,first_name,last_name,warehouse_id) VALUES (?,?,?,?)"

	employee := domain.Employee{
		CardNumberID: "12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	t.Run("Save OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, employee)
		assert.NoError(t, err)
		assert.Equal(t, 1, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, employee)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error WarehouseNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1452, Message: "warehouse"})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, employee)
		assert.Error(t, err)
		assert.Equal(t, ErrWarehouseNotfound, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error Result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnError(&mysql.MySQLError{Number: 1000, Message: ""})

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, employee)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error LastInserId", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		repo := NewRepository(db)
		lastId, err := repo.Save(ctx, employee)
		assert.Error(t, err)
		assert.Equal(t, 0, lastId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_Update(t *testing.T) {
	ctx := context.Background()

	query := "UPDATE employees SET first_name=?, last_name=?, warehouse_id=?  WHERE id=?"

	employee := domain.Employee{
		ID:           1,
		CardNumberID: "12",
		FirstName:    "Juan",
		LastName:     "Perez",
		WarehouseID:  1,
	}

	t.Run("Update OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(db)
		err = repo.Update(ctx, employee)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)
		err = repo.Update(ctx, employee)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error WarehouseNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1452})

		repo := NewRepository(db)
		err = repo.Update(ctx, employee)
		assert.Error(t, err)
		assert.Equal(t, ErrWarehouseNotfound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnError(&mysql.MySQLError{})

		repo := NewRepository(db)
		err = repo.Update(ctx, employee)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update Error Result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		repo := NewRepository(db)
		err = repo.Update(ctx, employee)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_Delete(t *testing.T) {
	ctx := context.Background()

	query := "DELETE FROM employees WHERE id=?"

	t.Run("Delete OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectExec().
			WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewRepository(db)
		err = repo.Delete(ctx, 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)
		err = repo.Delete(ctx, 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Exec", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(1).WillReturnError(sql.ErrNoRows)

		repo := NewRepository(db)
		err = repo.Delete(ctx, 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error Rows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(ErrDatabase))

		repo := NewRepository(db)
		err = repo.Delete(ctx, 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete Error NotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(query)).
			ExpectExec().WithArgs(1).WillReturnResult(driver.RowsAffected(0))

		repo := NewRepository(db)
		err = repo.Delete(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_GetAllInboundOrders(t *testing.T) {
	ctx := context.Background()

	query := "SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(i.id) FROM employees e LEFT JOIN inbound_orders i ON e.id = i.employee_id GROUP BY e.id;"

	data := []domain.EmployeeWithInboundOrders{
		{
			ID:                 1,
			FirstName:          "Juan",
			LastName:           "Perez",
			WarehouseID:        1,
			CardNumberID:       "12",
			InboundOrdersCount: 2,
		},
		{
			ID:                 2,
			FirstName:          "Carlos",
			LastName:           "Perez",
			WarehouseID:        1,
			CardNumberID:       "121",
			InboundOrdersCount: 1,
		},
	}

	t.Run("GetAllInoudnOrders OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "COUNT(i.id)"})
		for _, d := range data {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName, d.WarehouseID, d.InboundOrdersCount)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := NewRepository(db)

		employees, err := repo.GetAllInoundOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, data, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllInoundOrders Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("Error data base"))

		repo := NewRepository(db)

		employees, err := repo.GetAllInoundOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAllInoundOrders Error Scan", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		for _, d := range data {
			rows.AddRow(d.ID, d.CardNumberID, d.FirstName, d.LastName, d.WarehouseID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

		repo := NewRepository(db)

		employees, err := repo.GetAllInoundOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func Test_Repository_GetWithInoundOrders(t *testing.T) {
	ctx := context.Background()

	query := "SELECT e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id, COUNT(i.id) FROM employees e LEFT JOIN inbound_orders i ON e.id = i.employee_id WHERE e.id=? GROUP BY e.id;"

	data := domain.EmployeeWithInboundOrders{
		ID:                 1,
		FirstName:          "Juan",
		LastName:           "Perez",
		WarehouseID:        1,
		CardNumberID:       "12",
		InboundOrdersCount: 2,
	}

	t.Run("GetWithInoundOrders OK", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id", "COUNT(i.id)"})
		rows.AddRow(data.ID, data.CardNumberID, data.FirstName, data.LastName, data.WarehouseID, data.InboundOrdersCount)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

		repo := NewRepository(db)

		employee, err := repo.GetWithInboundOrder(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, data, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetWithInoudOrders Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(sql.ErrConnDone)

		repo := NewRepository(db)

		employee, err := repo.GetWithInboundOrder(ctx, 1)
		assert.Error(t, err)
		assert.Empty(t, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
