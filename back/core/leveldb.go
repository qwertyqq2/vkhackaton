package core

import (
	"database/sql"
	"os"

	"github.com/qwertyqq2/filebc/core/types"
)

const (
	DbName = "blockchain.db"
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

func (l *levelDB) insertBlock(hash, block string) error {
	_, err := l.db.Exec("INSERT INTO Blockchain (Hash, Block) VALUES ($1, $2);",
		hash,
		block,
	)
	return err
}

func (l *levelDB) getBlocks() (types.Blocks, error) {
	rows, err := l.db.Query("Select Block from Blockchain")
	if err != nil {
		return nil, err
	}
	blocksstr := make([]string, 0)
	var bs string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&bs)
		blocksstr = append(blocksstr, bs)
	}
	blocks := make([]*types.Block, 0)
	for _, bs := range blocksstr {
		b, err := types.DeserializeBlock(bs)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}
