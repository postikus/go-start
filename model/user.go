package model

//go:generate easytags $GOFILE
type User struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Balance *Balance `json:"balance"`
}
