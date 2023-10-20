package main

import (
	todo "MaksZhukGO"
	"log"
)

func main() {
	srv := new(todo.Server)
	if err := srv.Run("8000"); err != nil {
		log.Fatal("error occurred while running http server: %s", err.Error())
	}
}
