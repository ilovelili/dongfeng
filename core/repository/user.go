package repository

import (
	"github.com/ilovelili/dongfeng/core/model"
)

// User user repository
type User struct{}

// NewUserRepository new user repository
func NewUserRepository() *User {
	db().AutoMigrate(&model.User{})
	return new(User)
}

// Save save user
func (r *User) Save(user *model.User) error {
	return db().Save(user).Error
}

// Find find
func (r *User) Find(id uint) (*model.User, error) {
	user := new(model.User)
	err := db().First(user, id).Error
	return user, err
}

// FindByEmail find by email
func (r *User) FindByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := db().Where("email = ?", email).First(user).Error
	return user, err
}
