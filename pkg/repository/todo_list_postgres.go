package repository

import (
	"fmt"
	"github.com/kkosim/todo-app"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"strings"
)

type TodoListPostgres struct {
	db *gorm.DB
}

func NewTodoListPostgres(db *gorm.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx := r.db.Begin()
	var id int
	createListQuery := fmt.Sprintf("insert into %s (title, description) values ($1, $2) returning id", todoListsTable)

	err := tx.Raw(createListQuery, list.Title, list.Description).Scan(&id).Error
	if err != nil {
		logrus.Error("couldn't create new todoList")
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("insert into %s (user_id, list_id) values ($1, $2)", userListsTable)
	err = tx.Exec(createUserListQuery, userId, id).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	return id, nil
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on "+
		"tl.id=ul.list_id where ul.user_id=?", todoListsTable, userListsTable)
	err := r.db.Raw(query, userId).Scan(&lists).Error

	return lists, err
}

func (r *TodoListPostgres) GetById(userId int, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl inner join %s ul on "+
		"tl.id=ul.list_id where ul.user_id=? and ul.list_id=?", todoListsTable, userListsTable)
	err := r.db.Raw(query, userId, listId).Scan(&list).Error

	return list, err
}

func (r *TodoListPostgres) Delete(userId int, listId int) error {
	query := fmt.Sprintf("delete from %s tl using %s ul "+
		"where tl.id=ul.list_id and ul.user_id=? and ul.list_id=?", todoListsTable, userListsTable)
	err := r.db.Exec(query, userId, listId).Error

	return err
}

func (r *TodoListPostgres) Update(userId int, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	//title = $1
	//description = $1
	//title = $1 description = $2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("update %s tl set %s from %s ul "+
		"where tl.id=ul.list_id and ul.list_id=$%d and ul.user_id=$%d",
		todoListsTable, setQuery, userListsTable, argId, argId+1)
	args = append(args, listId, userId)

	log.Println("updateQuery: ", query)
	log.Println("args: ", args)

	err := r.db.Exec(query, args...).Error
	//err.Save()
	return err
}
