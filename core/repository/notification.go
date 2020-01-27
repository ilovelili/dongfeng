package repository

import "github.com/ilovelili/dongfeng/core/model"

// Notification notification repository
type Notification struct{}

// NewNotificationRepository new notification repository
func NewNotificationRepository() *Notification {
	db().AutoMigrate(&model.Notification{})
	return new(Notification)
}

// Save save notification
func (r *Notification) Save(notification *model.Notification) error {
	return db().Save(notification).Error
}

// SetRead set notification read
func (r *Notification) SetRead(notificationIDs []int) error {
	return db().Model(&model.Notification{}).Where("id IN (?)", notificationIDs).Update("read", true).Error
}

// FindByEmail find by email
func (r *Notification) FindByEmail(email string) ([]*model.Notification, error) {
	notifications := []*model.Notification{}
	err := db().Where("user = ? AND `read` = 0", email).Find(&notifications).Error
	return notifications, err
}
