package product_records

import (
	"context"
	"testing"

	"github.com/mercadolibre/fury_bootcamp-go-w7-s4-8-3/internal/domain"
	"github.com/stretchr/testify/assert"
)

type dummyRepo struct {
	Exist   bool
	ValidID bool
	ID      int
	Err     error
}

func (d dummyRepo) ValidateProductID(ctx context.Context, pid int) bool {
	return d.ValidID
}
func (d dummyRepo) Store(ctx context.Context, pr domain.ProductRecord) (int, error) {
	return d.ID, d.Err
}

var dummyProductRecordArg = domain.ProductRecord{
	LastUpdateDate: "12-12-2022",
	PurchasePrice:  12.34,
	SalePrice:      21.21,
	ProductID:      12,
}

var dummyProductRecordRes = domain.ProductRecord{
	ID:             14,
	LastUpdateDate: "12-12-2022",
	PurchasePrice:  12.34,
	SalePrice:      21.21,
	ProductID:      12,
}

func createTestService(d dummyRepo) (Service, context.Context) {
	return &service{r: d}, context.Background()
}

func TestValidateProductID_true(t *testing.T) {
	s, c := createTestService(dummyRepo{
		ValidID: true,
	})
	// should return true
	valid := s.ValidateProductID(c, 12)

	assert.True(t, valid)
}

func TestValidateProductID_false(t *testing.T) {
	s, c := createTestService(dummyRepo{
		ValidID: false,
	})
	// should return false
	valid := s.ValidateProductID(c, 12)

	assert.False(t, valid)
}

func TestCreate(t *testing.T) {
	s, c := createTestService(dummyRepo{
		ID:  14,
		Err: nil,
	})
	// should return dummyProductRecordRes, nil
	pr, err := s.Create(c, dummyProductRecordArg)

	assert.NoError(t, err)
	assert.Equal(t, dummyProductRecordRes, pr)
}

func TestCreate_ErrDatabase(t *testing.T) {
	s, c := createTestService(dummyRepo{
		ID:  14,
		Err: ErrDatabase,
	})
	// should return domain.ProductRecord{}, ErrDatabse
	pr, err := s.Create(c, dummyProductRecordArg)

	assert.Error(t, err)
	assert.Equal(t, ErrDatabase, err)
	assert.Equal(t, domain.ProductRecord{}, pr)
}
