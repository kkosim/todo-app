package main

import (
	"github.com/kkosim/todo-app"
	"github.com/kkosim/todo-app/pkg/handler"
	"github.com/kkosim/todo-app/pkg/repository"
	"github.com/kkosim/todo-app/pkg/service"
	"log"
)

func main() {

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(todo.Server)
	err := srv.Run("8000", handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
