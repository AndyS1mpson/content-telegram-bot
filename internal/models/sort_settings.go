package models

const (
	// SortOrderAsc Сортировка по возрастанию
	SortOrderAsc SortOrder = "ASC"

	// SortOrderDesc Сортировка по убыванию.
	SortOrderDesc SortOrder = "DESC"
)

// SortOrder тип для обозначения направленности сортировки
type SortOrder string

// SortSetting настройка сортировки для произвольного поля
type SortSetting struct {
	SortField string
	SortOrder SortOrder
}
