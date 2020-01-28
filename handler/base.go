package handler

import (
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
	"github.com/ilovelili/dongfeng/util"
)

var (
	config           = util.LoadConfig()
	userRepo         = repository.NewUserRepository()
	classRepo        = repository.NewClassRepository()
	pupilRepo        = repository.NewPupilRepository()
	teacherRepo      = repository.NewTeacherRepository()
	notificationRepo = repository.NewNotificationRepository()
)

// notify notification
func notify(notification *model.Notification) {
	go func() {
		notificationRepo.Save(notification)
	}()
}
