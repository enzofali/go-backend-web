package domain


type InboundOrder struct {
	ID 				int `json:"id"`
	OrderDate 		string `json:"order_date"`
	OrderNumber 	string `json:"order_number"`
	EmployeeID 		int `json:"employee_id"`
	ProductBatchID  int	`json:"product_batch_id"`
	WarehouseID 	int `json:"warehouse_id"`
}

type InboundOrderRequest struct {
	OrderDate 		string `json:"order_date" validate:"required"`
	OrderNumber 	string `json:"order_number" validate:"required"`
	EmployeeID 		int `json:"employee_id" validate:"required"`
	ProductBatchID  int	`json:"product_batch_id" validate:"required"`
	WarehouseID 	int `json:"warehouse_id" validate:"required"`
}