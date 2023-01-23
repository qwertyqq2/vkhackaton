package core

import (
	"database/sql"
	"os"
)

const (
	DbName = "files.db"
)

type levelDB struct {
	db *sql.DB
}

func NewLevelDB() (*levelDB, error) {
	fileDb, err := os.Create(DbName)
	if err != nil {
		return nil, err
	}
	fileDb.Close()
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE Blockchain (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Hash VARCHAR(44) UNIQUE,
		Block TEXT
	);
	`)
	if err != nil {
		return nil, err
	}
	return &levelDB{db: db}, nil
}

func loadLevel() (*levelDB, error) {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return nil, err
	}
	return &levelDB{
		db: db,
	}, nil
}
 
