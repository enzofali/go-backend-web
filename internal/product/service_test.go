package product

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

// define a dummy repository struct that implements the Repository interface for testing purposes
type stubRepo struct {
	Id          int
	Err         error
	Products    []domain.Product
	P           domain.Product
	Exist       bool
	Records     int
	Reports     []domain.Report
	Description string
	ValidID     bool
}

// All methods simply return the struct's initially defined members
func (d stubRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	return d.Products, d.Err
}
func (d stubRepo) Get(ctx context.Context, id int) (domain.Product, error) {
	return d.P, d.Err
}
func (d stubRepo) Exists(ctx context.Context, productCode string) bool {
	return d.Exist
}
func (d stubRepo) Save(ctx context.Context, p domain.Product) (int, error) {
	return d.Id, d.Err
}
func (d stubRepo) Update(ctx context.Context, p domain.Product) error {
	return d.Err
}
func (d stubRepo) Delete(ctx context.Context, id int) error {
	return d.Err
}
func (d stubRepo) ValidateProductID(ctx context.Context, pid int) bool {
	return d.ValidID
}
func (d stubRepo) GetOneReport(ctx context.Context, id int) (int, string, error) {
	return d.Records, d.Description, d.Err
}
func (d stubRepo) GetAllReports(ctx context.Context) ([]domain.Report, error) {
	return d.Reports, d.Err
}
func (d stubRepo) StoreType(ctx context.Context, name string) (int, error) {
	return d.Id, d.Err
}

// dummy product
var dummyP = domain.Product{
	ID:             13,
	Description:    "pepe",
	ExpirationRate: 4,
	FreezingRate:   5,
	Height:         12.3,
	Length:         11.9,
	Netweight:      10.9,
	ProductCode:    "alala",
	RecomFreezTemp: 123.0,
	Width:          4.2,
	ProductTypeID:  100,
	SellerID:       1,
}

// product with different product code, for testing Update method
var updateP = domain.Product{
	ID:             13,
	Description:    "pepe",
	ExpirationRate: 4,
	FreezingRate:   5,
	Height:         12.3,
	Length:         11.9,
	Netweight:      10.9,
	ProductCode:    "exists",
	RecomFreezTemp: 123.0,
	Width:          4.2,
	ProductTypeID:  100,
	SellerID:       1,
}

// product with no pre-defined id, for testing the Create method
var dummyP_NoID = domain.Product{
	Description:    "pepe",
	ExpirationRate: 4,
	FreezingRate:   5,
	Height:         12.3,
	Length:         11.9,
	Netweight:      10.9,
	ProductCode:    "alala",
	RecomFreezTemp: 123.0,
	Width:          4.2,
	ProductTypeID:  100,
	SellerID:       1,
}

// dummy product inventory, for testing GetAll
var dummyProducts = []domain.Product{
	dummyP,
}

var dummyReport = []domain.Report{
	{Count: 1, Description: "test", ProductID: 1},
}

var dummyReports = []domain.Report{
	{Count: 14, Description: "pepe", ProductID: 12},
	{Count: 15, Description: "not pepe", ProductID: 21},
}

// returns a service using the dummy repo and a default context, unused but necessary to comply with method signatures
func createTestService(d stubRepo) (Service, context.Context) {
	return &service{r: d}, context.Background()
}

func TestGetAll(t *testing.T) {
	s, c := createTestService(stubRepo{
		Products: dummyProducts,
		Err:      nil,
	})
	// should return dummyProducts, nil
	products, err := s.GetAll(c)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, dummyProducts, products)
}

func TestGetAll_Err(t *testing.T) {
	s, c := createTestService(stubRepo{
		Products: dummyProducts,
		Err:      ErrDatabase,
	})
	// should return []domain.Product{}, ErrDatabase
	products, err := s.GetAll(c)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, 0, len(products))
	assert.Equal(t, []domain.Product{}, products)
}

func TestGetByID(t *testing.T) {
	s, c := createTestService(stubRepo{
		P:   dummyP,
		Err: nil,
	})
	// should return dummyP, nil
	product, err := s.GetByID(c, 13)

	assert.NoError(t, err)
	assert.Equal(t, dummyP, product)
}

func TestGetByID_Err(t *testing.T) {
	s, c := createTestService(stubRepo{
		P:   dummyP,
		Err: ErrNotFound,
	})
	// should return domain.Product{}, nil
	product, err := s.GetByID(c, 13)

	assert.Error(t, err)
	assert.Equal(t, domain.Product{}, product)
	assert.Equal(t, ErrNotFound, err)
}

func TestCreate(t *testing.T) {
	s, c := createTestService(stubRepo{
		Id:       13,
		Products: []domain.Product{},
		P:        domain.Product{},
		Err:      nil,
		Exist:    false,
	})
	// should return dummyP, nil
	product, err := s.Create(c, dummyP_NoID)

	assert.NoError(t, err)
	assert.Equal(t, dummyP, product)
}

func TestCreate_ErrExists(t *testing.T) {
	s, c := createTestService(stubRepo{
		Id:       13,
		Products: []domain.Product{},
		P:        domain.Product{},
		Err:      nil,
		Exist:    true,
	})
	// should return domain.Product{}, ErrExists
	product, err := s.Create(c, dummyP_NoID)

	assert.Error(t, err)
	assert.Equal(t, ErrExists, err)
	assert.Equal(t, domain.Product{}, product)
}

func TestCreate_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		Id:       13,
		Products: []domain.Product{},
		P:        domain.Product{},
		Err:      ErrDatabase,
		Exist:    false,
	})
	// should return domain.Product{}, nil
	product, err := s.Create(c, dummyP_NoID)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, domain.Product{}, product)
}

func TestUpdate(t *testing.T) {
	s, c := createTestService(stubRepo{
		P:        dummyP,
		Products: dummyProducts,
		Err:      nil,
		Exist:    false,
	})
	// should return dummyP, nil
	p, err := s.Update(c, dummyP)

	assert.NoError(t, err)
	assert.Equal(t, dummyP, p)
}

func TestUpdate_ErrExists(t *testing.T) {
	s, c := createTestService(stubRepo{
		P:     dummyP,
		Err:   nil,
		Exist: true,
	})
	// should return domain.Product{}, ErrExists
	p, err := s.Update(c, updateP)

	assert.Error(t, err)
	assert.Equal(t, ErrExists, err)
	assert.Equal(t, domain.Product{}, p)
}

func TestUpdate_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		P:     dummyP,
		Err:   ErrDatabase,
		Exist: false,
	})
	// should return domain.Product{}, ErrDatabase
	p, err := s.Update(c, updateP)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, domain.Product{}, p)
}

func TestDelete(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err: nil,
	})
	// should return nil
	err := s.Delete(c, 13)

	assert.NoError(t, err)
}

func TestDelete_ErrNotFound(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err: ErrNotFound,
	})
	// should return ErrNotFound
	err := s.Delete(c, 13)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestDelete_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err: ErrDatabase,
	})
	// should return ErrDatabase
	err := s.Delete(c, 13)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
}

func TestValidateProductID_true(t *testing.T) {
	s, c := createTestService(stubRepo{
		ValidID: true,
	})
	// should return true
	valid := s.ValidateProductID(c, 1)

	assert.Equal(t, true, valid)
}

func TestValidateProductID_false(t *testing.T) {
	s, c := createTestService(stubRepo{
		ValidID: false,
	})
	// should return true
	valid := s.ValidateProductID(c, 1)

	assert.Equal(t, false, valid)
}

func TestGetOneReport(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err:         nil,
		Records:     1,
		Description: "test",
	})
	// should return 1, "test", nil
	report, err := s.GetOneReport(c, 1)

	assert.NoError(t, err)
	assert.Equal(t, dummyReport, report)
}

func TestGetOneReport_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err:     ErrDatabase,
		Reports: dummyReport,
	})
	// should return 0, "", ErrDatabase
	report, err := s.GetOneReport(c, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, []domain.Report{}, report)
}

func TestGetAllReports(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err:     nil,
		Reports: dummyReports,
	})
	// should return 14, nil
	reports, err := s.GetAllReports(c)

	assert.NoError(t, err)
	assert.Equal(t, dummyReports, reports)
}

func TestGetAllReports_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		Err:     ErrDatabase,
		Reports: dummyReports,
	})
	// should return 0, ErrDatabase
	reports, err := s.GetAllReports(c)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, []domain.Report{}, reports)
}

func TestCreateType_Ok(t *testing.T) {
	s, c := createTestService(stubRepo{
		Id:  1,
		Err: nil,
	})
	// should return 1, nil
	typeId, err := s.CreateType(c, "pepe")

	assert.NoError(t, err)
	assert.Equal(t, 1, typeId)
}

func TestCreateType_ErrDatabase(t *testing.T) {
	s, c := createTestService(stubRepo{
		Id:  1,
		Err: ErrDatabase,
	})
	// should return 0, ErrDatabase
	typeId, err := s.CreateType(c, "pepe")

	assert.Error(t, err)
	assert.Equal(t, 0, typeId)
	assert.Equal(t, ErrDatabase, err)
}
