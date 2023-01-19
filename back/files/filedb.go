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
	CREATE TABLE Users (
		Address varchar,
		Balance integer
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

type wrapper struct {
	Addr string `json:"Address"`
	Bal  int    `json:"Balance"`
}

func (l *levelDB) newUser(address string) error {
	_, err := l.db.Exec("INSERT INTO Users (Address, Balance) VALUES ($1, $2);",
		address,
		int(0),
	)
	return err
}

func (l *levelDB) existUser(address string) bool {
	row := l.db.QueryRow("Select Balance From Users Where Address=$1", address)
	var res interface{}
	row.Scan(&res)
	if res == nil {
		return false
	}
	return true
}

func (l *levelDB) addBalance(address string, delta uint64) error {
	if l.existUser(address) {
		_, err := l.db.Exec("Update Users Set Balance=$1 Where Address=$2;",
			int(delta),
			address,
		)
		return err
	}
	err := l.newUser(address)
	if err != nil {
		return err
	}
	_, err = l.db.Exec("Update Users Set Balance=$1 Where Address=$2;",
		int(delta),
		address,
	)
	return err
}

func (l *levelDB) getBalance(address string) (uint64, error) {
	row := l.db.QueryRow("Select Balance From Users Where Address=$1", address)
	var user wrapper
	err := row.Scan(&user.Bal)
	if err != nil {
		return 0, err
	}
	return uint64(user.Bal), nil
}

func (l *levelDB) getUsers() ([]wrapper, error) {
	rows, err := l.db.Query("Select * from Users")
	if err != nil {
		return nil, err
	}
	addrs := make([]wrapper, 0)
	var as wrapper
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&as.Addr, &as.Bal)
		addrs = append(addrs, as)
	}
	return addrs, nil
}
