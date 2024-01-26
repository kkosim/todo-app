package repository

import (
	"github.com/kkosim/todo-app"
	"gorm.io/gorm"
)

type TodoItemPostgres struct {
	db *gorm.DB
}

func NewTodoItemPostgres(db *gorm.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx := r.db.Begin()
	var listIt todo.ListItem
	err := tx.Table(todoItemsTable).Create(&item).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	listIt = todo.ListItem{
		ListId: listId,
		ItemId: item.Id,
	}
	err = tx.Table(listsItemsTable).Create(&listIt).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit().Error
}
