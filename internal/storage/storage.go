package storage

import (
	"github.com/MdZunaed/students-api/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	UpdateStudent(name string, email string, age int, id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	GetStudentById(id int64) (types.Student, error)
	DeleteStudentById(id int64) (error)
}
