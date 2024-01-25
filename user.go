package todo

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" gorm:"column:name" binding:"required"`
	Username string `json:"username" gorm:"column:username" binding:"required"`
	Password string `json:"password" gorm:"column:password_hash" binding:"required"`
}
