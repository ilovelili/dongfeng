package repository

import "github.com/ilovelili/dongfeng/core/model"

// Notification notification repository
type Notification struct{}

// NewNotificationRepository new notification repository
func NewNotificationRepository() *Notification {
	db().AutoMigrate(&model.Notification{})
	return new(Notification)
}

// Save save Notification
func (r *Notification) Save(Notification *model.Notification) error {
	return db().Save(Notification).Error
}

// FindByEmail find by email
func (r *Notification) FindByEmail(email string) ([]*model.Notification, error) {
	notifications := []*model.Notification{}
	err := db().Where("user = ?", email).Find(&notifications).Error
	return notifications, err
}
