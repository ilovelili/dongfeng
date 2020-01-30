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
		return util.ResponseError(c, "500-101", "failed to get notifications", err)
	}

	return c.JSON(http.StatusOK, notifications)
}

// SetNotificationsRead POST /notifications
func SetNotificationsRead(c echo.Context) error {
	type setReadReq struct {
		IDs []int `json:"ids"`
	}
	req := new(setReadReq)
	if err := c.Bind(req); err != nil {
		return util.ResponseError(c, "400-103", "failed to bind notifications", err)
	}

	if err := notificationRepo.SetRead(req.IDs); err != nil {
		return util.ResponseError(c, "500-102", "failed to save notifications", err)
	}

	return c.NoContent(http.StatusOK)
}
