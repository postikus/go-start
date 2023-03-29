package model

type Balance struct {
	UserID int64   `json:"user_id"`
	Amount float32 `json:"amount"`
}
