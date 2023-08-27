package repo

import (
	. "lightweightpipline/internal/biz/aggregates"
)

type OrderRepo struct {
	BaseRepo
}

func (o *OrderRepo) GetAllOrders() ([]Command, error) {
	dto := []Command{}
	result := o.Context.GetDb().Find(&dto)
	return dto, result.Error
}
