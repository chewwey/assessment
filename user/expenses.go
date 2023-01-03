package user

import "database/sql"

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float32  `json:"amount"`
	Note   string   `json:"note"`
	Tag    []string `json:"tags"`
}

type Err struct {
	Message string
}

type Handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *Handler {
	return &Handler{db}
}
