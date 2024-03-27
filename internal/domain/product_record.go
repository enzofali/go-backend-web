package domain

type ProductRecord struct {
	ID             int     `json:"id"`
	LastUpdateDate string  `json:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price" validate:"required"`
	SalePrice      float64 `json:"sale_price" validate:"required"`
	ProductID      int     `json:"product_id" validate:"required"`
}

// ProductRecordRequest exists solely for Swaggo/Swagger documentation purposes
type ProductRecordRequest struct {
	LastUpdateDate string  `json:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price" validate:"required"`
	SalePrice      float64 `json:"sale_price" validate:"required"`
	ProductID      int     `json:"product_id" validate:"required"`
}
