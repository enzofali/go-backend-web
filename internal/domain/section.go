package domain

type Section struct {
	ID                 int `json:"section_id"`
	SectionNumber      int `json:"section_number" validate:"required"`
	CurrentTemperature int `json:"current_temperature" validate:"required"`
	MinimumTemperature int `json:"minimum_temperature" validate:"required"`
	CurrentCapacity    int `json:"current_capacity" validate:"required"`
	MinimumCapacity    int `json:"minimum_capacity" validate:"required"`
	MaximumCapacity    int `json:"maximum_capacity" validate:"required"`
	WarehouseID        int `json:"warehouse_id" validate:"required"`
	ProductTypeID      int `json:"product_type_id" validate:"required"`
}

type SectionReportProducts struct {
	//Section
	ID            int `json:"section_id"`
	SectionNumber int `json:"section_number"`
	ProductCount  int `json:"product_count"`
}
