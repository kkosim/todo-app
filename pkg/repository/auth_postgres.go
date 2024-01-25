package repository

import (
	"fmt"
	"github.com/kkosim/todo-app"
	"github.com/sirupsen/logrus"
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

	err := r.db.Raw(query, user.Name, user.Username, user.Password).Scan(&id).Error
	if err != nil {
		logrus.Error("couldn't create new user")
	}

	//err := r.db.Table(userTable).Create(userAuth{
	//	name:     user.Name,
	//	username: user.Username,
	//	password: user.Password,
	//}).Error
	//
	//if err == nil {
	//	id = user.Id
	//}
	//if err != nil {
	//	logrus.Error("couldn't create new user")
	//}
	//
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User

	err := r.db.Table(userTable).Where("username=? and password_hash=?", username, password).Scan(&user).Error
	if err != nil {
		return todo.User{}, err
	}
	//fmt.Println("user:", user)
	return user, nil
}
