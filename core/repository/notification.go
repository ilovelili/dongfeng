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

// SaveAll save all notifications
func (r *Notification) SaveAll(notifications []*model.Notification) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := tx.Save(notification).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// FindByEmail find by email
func (r *Notification) FindByEmail(email string) ([]*model.Notification, error) {
	notifications := []*model.Notification{}
	err := db().Where("user = ?", email).Find(&notifications).Error
	return notifications, err
}
