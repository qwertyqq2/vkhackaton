package files

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/qwertyqq2/filebc/crypto"
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
	CREATE TABLE Files (
		Id varchar,
		Rand varchar,
		File text
	);
	`)
	if err != nil {
		return nil, err
	}
	return &levelDB{db: db}, nil
}

func LoadLevel() (*levelDB, error) {
	db, err := sql.Open("sqlite3", DbName)
	if err != nil {
		return nil, err
	}
	return &levelDB{
		db: db,
	}, nil
}

func (l *levelDB) insertFile(f *File) error {
	fstr, err := f.SerializeFile()
	if err != nil {
		return err
	}
	_, err = l.db.Exec("INSERT INTO Files VALUES ($1, $2, $3)",
		crypto.Base64EncodeString(f.Id),
		crypto.Base64EncodeString(f.rand),
		fstr,
	)
	if err != nil {
		return err
	}
	return nil
}

func (l *levelDB) allFiles() ([]*File, error) {
	fsarr, err := l.getFiles()
	if err != nil {
		return nil, err
	}
	files := make([]*File, 0)
	for _, fs := range fsarr {
		f, err := DeserializeFile(fs)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (l *levelDB) getFiles() ([]string, error) {
	rows, err := l.db.Query("Select File from Files")
	if err != nil {
		return nil, err
	}
	filesStr := make([]string, 0)
	var fs string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&fs)
		filesStr = append(filesStr, fs)
	}
	return filesStr, nil

}
