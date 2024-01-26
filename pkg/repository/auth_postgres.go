package repository

import (
	"github.com/kkosim/todo-app"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

//type userAuth struct {
//	name     string
//	username string
//	password string
//}

func newAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	err := r.db.Table(userTable).Create(&user).Error
	if err != nil {
		logrus.Error("couldn't create new user: ", err)
	}
	return user.Id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User

	err := r.db.Table(userTable).Where("username=? and password_hash=?", username, password).Scan(&user).Error
	if err != nil {
		return todo.User{}, err
	}
	return user, nil
}
