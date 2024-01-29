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

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// list doesn't exist or belongs to another user
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId int, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId int, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId int, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
