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
	query := fmt.Sprintf("insert into %s (name, username, password_hash) values ($1, $2, $3) returning id", userTable)

	row := r.db.Raw(query, user.Name, user.Username, user.Password)
	row.Find(&id)
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	err := r.db.Where("username=$1 and password_hash=$2", username, password).Find(&user)
	if err != nil {
		return todo.User{}, err.Error
	}
	return user, nil
}
