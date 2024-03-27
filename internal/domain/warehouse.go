package domain

type Warehouse struct {
	ID                 int    `json:"id"`
	Address            string `json:"address" validate:"required"`
	Telephone          string `json:"telephone" validate:"required"`
	WarehouseCode      string `json:"warehouse_code" validate:"required"`
	MinimumCapacity    int    `json:"minimum_capacity" validate:"required"`
	MinimumTemperature int    `json:"minimum_temperature" validate:"required"`
}

// create struct for validate field empty
type WarehouseRequest struct {
	Address            string `json:"address" validate:"required"`
	Telephone          string `json:"telephone" validate:"required"`
	WarehouseCode      string `json:"warehouse_code" validate:"required"`
	MinimumCapacity    int    `json:"minimum_capacity" validate:"required"`
	MinimumTemperature int    `json:"minimum_temperature" validate:"required"`
}
