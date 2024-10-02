package sql

import "content-telegram-bot/internal/models"


const (
	// OrderDirectionAsc сортировка по возрастанию
	OrderDirectionAsc OrderDirection = "ASC"

	// OrderDirectionDesc сортировка по убыванию
	OrderDirectionDesc OrderDirection = "DESC"
)

// SortOrderToDirection маппинг порядка сортировки
var SortOrderToDirection = map[models.SortOrder]OrderDirection{
	models.SortOrderAsc:  OrderDirectionAsc,
	models.SortOrderDesc: OrderDirectionDesc,
}

// OrderDirection направление сортировки для колонки
type OrderDirection string

// OrderColumn колонка, доступная для сортировки
type OrderColumn string

// OrderByColumn тип сортировки для колонки
type OrderByColumn struct {
	Column    OrderColumn
	Direction OrderDirection
}

// OrderBy список сортировок для запроса
type OrderBy struct {
	Columns []OrderByColumn
}
