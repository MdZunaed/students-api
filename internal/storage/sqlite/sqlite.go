package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MdZunaed/students-api/internal/config"
	"github.com/MdZunaed/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func InitSqlite(cfg *config.Config) (*Sqlite, error) {
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

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	//stmt,err:= s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	//stmt,err:= s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
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

// / Safe update, will update single field too
func (s *Sqlite) UpdateStudent(name *string, email *string, age *int, id int64) (*types.Student, error) {
	var student types.Student
	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
		Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	log.Printf("student after fetch, name:%s, email:%s, age:%d", student.Name, student.Email, student.Age)
	if err != nil {
		return nil, fmt.Errorf("not found")
	}
	// Override only the non-nil fields
	if name != nil {
		student.Name = *name
	}
	if email != nil {
		student.Email = *email
	}
	if age != nil {
		student.Age = *age
	}
	stmt, err := s.Db.Prepare(`
		UPDATE students
	 	SET name = ?, email = ?, age = ?
	 	WHERE id = ?
	`)
	if err != nil {
		log.Printf("Prepare error: %v\n", err)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(student.Name, student.Email, student.Age, id)
	if err != nil {
		log.Printf("Exec error: %v\n", err)
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &student, nil
}

/// Update whole data, will wipe data if full data not given
// func (s *Sqlite) UpdateStudent(name string, email string, age int, id int64) (types.Student, error) {
// 	stmt, err := s.Db.Prepare(`
// 		UPDATE students
// 	 	SET name = ?, email = ?, age = ?
// 	 	WHERE id = ?
// 	`)
// 	if err != nil {
// 		log.Printf("Prepare error: %v\n", err)
// 		return types.Student{}, err
// 	}
// 	defer stmt.Close()
// 	result, err := stmt.Exec(name, email, age, id)
// 	if err != nil {
// 		log.Printf("Exec error: %v\n", err)
// 		return types.Student{}, err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return types.Student{}, err
// 	}
// 	if rowsAffected == 0 {
// 		return types.Student{}, fmt.Errorf("not found")
// 	}
// 	var student types.Student
// 	err = s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
// 		Scan(&student.Id, &student.Name, &student.Email, &student.Age)
// 	if err != nil {
// 		return types.Student{}, err
// 	}
// 	return student, nil
// }

func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("delete affected rows:", affected)
	if affected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}
