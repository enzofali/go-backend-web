package domain

type Buyer struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerRequest struct {
	CardNumberID string `json:"card_number_id" validate:"required"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	ID           int
}

type ReportBuyersPurchases struct {
	ID                   int    `json:"id"`
	Card_Number_ID       string `json:"card_number_id"`
	First_Name           string `json:"first_name"`
	Last_Name            string `json:"last_name"`
	Purchase_Order_Count int    `json:"purchase_order_count"`
}
