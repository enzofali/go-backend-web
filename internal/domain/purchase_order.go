package domain

type Purchase_Orders struct {
	ID int 					`json:"id"`
	Order_number string 	`json:"order_number" validate:"required"`
	Order_date string 		`json:"order_date" validate:"required"`
	Tracking_code string 	`json:"tracking_code" validate:"required"`
	Buyer_id int 			`json:"buyer_id" validate:"required"`
	Product_record_id int 	`json:"product_record_id" validate:"required"`
	Order_Status_id int 	`json:"order_status_id" validate:"required"`
}
