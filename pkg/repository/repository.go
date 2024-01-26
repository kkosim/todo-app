package repository

import (
	"github.com/kkosim/todo-app"
	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId int, listId int) (todo.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(userID, listId int, item todo.TodoItem) (int, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
