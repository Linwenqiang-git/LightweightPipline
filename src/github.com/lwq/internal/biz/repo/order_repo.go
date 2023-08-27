package repo

import (
	. "lightweightpipline/internal/biz/aggregates"
)

type IOrderRepo interface {
	GetAllOrders() ([]Command, error)
}
