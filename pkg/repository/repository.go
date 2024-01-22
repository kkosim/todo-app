package repository

import (
	"github.com/kkosim/todo-app"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
	}
}
