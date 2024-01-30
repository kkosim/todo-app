package repository

import (
	"fmt"
	"github.com/kkosim/todo-app"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	err := tx.Table(todoListsTable).Create(&list).Error
	if err != nil {
		logrus.Error("couldn't create new todoList")
		return 0, err
	}
	id := list.Id

	ul := todo.UserList{
		UserId: userId,
		ListId: id,
	}
	err = tx.Table(userListsTable).Create(&ul).Error
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

func (r *TodoItemPostgres) Update(userId int, itemId int, input todo.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	//title = $1
	//description = $1
	//title = $1 description = $2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("update %s ti set %s from %s li, %s ul "+
		"where ti.id=li.item_id and li.list_id=ul.list_id and ul.user_id=$%d and ti.id = $%d",
		todoItemsTable, setQuery, listsItemsTable, userListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	//log.Println("updateQuery: ", query)
	//log.Println("args: ", args)

	err := r.db.Exec(query, args...).Error
	//err.Save()
	return err
}
