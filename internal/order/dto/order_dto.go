package dto

type CreateOrderRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
}

type OrderResponse struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}
