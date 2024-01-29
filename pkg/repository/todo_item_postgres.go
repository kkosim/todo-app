package repository

import (
	"fmt"
	"github.com/kkosim/todo-app"
	"gorm.io/gorm"
	"strings"
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
		//tx.Rollback()
		return 0, err
	}
	itemId := item.Id
	listIt = todo.ListItem{
		ListId: listId,
		ItemId: itemId,
	}
	err = tx.Table(listsItemsTable).Create(&listIt).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit().Error
}

func (r *TodoItemPostgres) GetAll(userId int, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	err := r.db.Table(todoItemsTable+" ti ").Select("ti.id, ti.title, ti.description, ti.done").
		Joins("inner join "+listsItemsTable+" li on li.item_id = ti.id").
		Joins("inner join "+userListsTable+" ul on ul.list_id = li.list_id").
		Where(" li.list_id = $1 and ul.user_id = $2", listId, userId).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetById(userId int, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	err := r.db.Table(todoItemsTable+" ti ").
		Joins("inner join "+listsItemsTable+" li on li.item_id = ti.id").
		Joins("inner join "+userListsTable+" ul on ul.list_id = li.list_id").
		Where(" ti.id = $1 and ul.user_id = $2", itemId, userId).Scan(&item).Error
	if err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	//var item todo.TodoItem

	//subQuery := r.db.
	//	Table(listsItemsTable).
	//	Select("li.list_id").
	//	Joins("JOIN "+userListsTable+" ul ON li.list_id = ul.list_id").
	//	Where("ul.user_id = ?", userId).
	//	Where("li.item_id = ?", itemId) //.SubQuery()
	//
	//err := r.db.
	//	Table(todoItemsTable).
	//	Joins("JOIN "+listsItemsTable+" li ON "+todoItemsTable+".id = li.item_id").
	//	Where(todoItemsTable+".id = ?", itemId).
	//	Where(todoItemsTable+".id IN (?)", subQuery).
	//	Delete(&item).
	//	Error

	//err := r.db.
	//	Table(todoItemsTable).
	//	Joins("JOIN "+listsItemsTable+" li ON "+todoItemsTable+".id = li.item_id").
	//	Joins("JOIN "+userListsTable+" ul ON li.list_id = ul.list_id").
	//	Where(todoItemsTable+".id = ?", itemId).
	//	Where("ul.user_id = ?", userId).
	//	Delete(&item).Error
	//return err
	//gorm гавно а не orm

	query := fmt.Sprintf(`DELETE FROM ` + todoItemsTable + ` ti USING ` + listsItemsTable + ` li, ` + userListsTable + ` ul
	WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND
	ul.user_id = $1 AND ti.id = $2`)
	err := r.db.Exec(query, userId, itemId).Error
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

	//log.Println("updateQuery: ", query)
	//log.Println("args: ", args)

	err := r.db.Exec(query, args...).Error
	//err.Save()
	return err
}
