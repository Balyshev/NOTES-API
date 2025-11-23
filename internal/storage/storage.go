package storage

import "database/sql"

//Storage содержит подключение к БД
type Storage struct {
	db *sql.DB
}

//создаёт новое подключение к БД
func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

//Закрывает подключение к БД
func (s *Storage) Close() error {
	return s.db.Close()
}

//проверяет подключение к БД
func (s *Storage) Ping() error {
	return s.db.Ping()
}
