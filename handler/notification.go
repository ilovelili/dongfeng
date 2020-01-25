package handler

import (
	"net/http"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// Notifications GET /notifications
func Notifications(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	notificationRepo := repository.NewNotificationRepository()
	notifications, err := notificationRepo.FindByEmail(userInfo.Email)
	if err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-101", "failed to get notifications", err)
	}

	return c.JSON(http.StatusOK, notifications)
}
