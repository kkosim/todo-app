package repository

import (
	"fmt"
	"github.com/kkosim/todo-app"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

type userAuth struct {
	name     string
	username string
	password string
}

func newAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (name, username, password_hash) values ($1, $2, $3) returning id")

	row := r.db.Raw(query, user.Name, user.Name, user.Password)
	row.Find(&id)
	return id, nil
}
