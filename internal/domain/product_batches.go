package domain

type ProductBatches struct {
	ID                 int    `json:"id"`
	BatchNumber        int    `json:"batch_number" validate:"required"`
	CurrentQuantity    int    `json:"current_quantity" validate:"required"`
	CurrentTemperature int    `json:"current_temperature" validate:"required"`
	DueDate            string `json:"due_date" validate:"required"`
	InitialQuantity    int    `json:"initial_quantity" validate:"required"`
	ManufacturingDate  string `json:"manufacturing_date" validate:"required"`
	ManufacturingHour  string `json:"manufacturing_hour" validate:"required"`
	MinumumTemperature int    `json:"minumum_temperature" validate:"required"`
	ProductID          int    `json:"product_id" validate:"required"`
	SectionID          int    `json:"section_id" validate:"required"`
}
