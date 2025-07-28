package sqlite

import (
	"database/sql"
	"log"

	"github.com/MdZunaed/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL
	)`)

	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// age column was missing somehow, that's why needed to update manually
	// _, err := s.Db.Exec(`ALTER TABLE students ADD COLUMN age INTEGER`)
	// if err != nil {
	// 	log.Println("ALTER TABLE error:", err)
	// }
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?,?,?)")
	if err != nil {
		log.Printf("Prepare error: %v\n", err)
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		log.Printf("Exec error: %v\n", err)
		return 0, err
	}
	lstId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lstId, nil
}
