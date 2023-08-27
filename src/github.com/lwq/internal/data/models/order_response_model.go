package models

type OrderResponse struct {
	Status  string
	Order   string
	Message string
}

const (
	SuccessStatus = "success"
	ErrorStatus   = "error"
)
