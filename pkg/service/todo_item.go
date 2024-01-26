package service

import (
	"github.com/kkosim/todo-app"
	"github.com/kkosim/todo-app/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userID, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userID, listId)
	if err != nil {
		// list doesn't exist or belongs to another user
		return 0, err
	}

	return s.repo.Create(userID, listId, item)
}
