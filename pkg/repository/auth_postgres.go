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

	query := fmt.Sprintf("select * from %s where username=? and password_hash=?", userTable)
	err := r.db.Raw(query, username, password).Scan(&user).Error
	if err != nil {
		return todo.User{}, err
	}
	//err := r.db.Table(userTable).Where("username =?", username).Scan(user).Error
	fmt.Println("user:", user)
	return user, nil
}
