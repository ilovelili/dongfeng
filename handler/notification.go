package handler

import (
	"net/http"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetNotifications GET /notifications
func GetNotifications(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	notifications, err := notificationRepo.FindByEmail(userInfo.Email)
	if err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-101", "failed to get notifications", err)
	}

	return c.JSON(http.StatusOK, notifications)
}

// SaveNotifications POST /notifications
func SaveNotifications(c echo.Context) error {
	notifications := []*model.Notification{}
	if err := c.Bind(notifications); err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-103", "failed to bind notifications", err)
	}

	if err := notificationRepo.SaveAll(notifications); err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-102", "failed to save notifications", err)
	}

	return c.NoContent(http.StatusOK)
}
