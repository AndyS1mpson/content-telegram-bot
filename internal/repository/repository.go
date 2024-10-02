package repository

var pinTableName = "pinterest_pin"

// Repository описывает паттерн Repository для взаимодействия с базой данных
type Repository struct {
	db  querier
}

// New конструктор для Repository
func New(db querier) *Repository {
	return &Repository{
		db:  db,
	}
}
