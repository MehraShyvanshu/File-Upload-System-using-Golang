package infrastructure

import (
	"database/sql"
	"project/internal/models"

	_ "modernc.org/sqlite"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo() *SQLiteRepo {
	db, err := sql.Open("sqlite", "file:files.db")
	if err != nil {
		panic("Failed to open database: " + err.Error())
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// create table
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id TEXT PRIMARY KEY,
		name TEXT,
		url TEXT
	)`

	_, err = db.Exec(query)
	if err != nil {
		panic("Failed to create table: " + err.Error())
	}

	return &SQLiteRepo{db: db}
}

func (s *SQLiteRepo) Save(file models.File) error {
	_, err := s.db.Exec(
		"INSERT INTO files(id, name, url) VALUES (?, ?, ?)",
		file.ID, file.Name, file.URL,
	)
	return err
}

func (s *SQLiteRepo) List() ([]models.File, error) {
	rows, err := s.db.Query("SELECT id, name, url FROM files")
	if err != nil {
		return []models.File{}, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(&f.ID, &f.Name, &f.URL); err != nil {
			return []models.File{}, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}
