package model

type Order struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}
