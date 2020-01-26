package handler

import (
	"net/http"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// Login POST /login
func Login(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	user, err := userRepo.FindByEmail(userInfo.Email)
	if err != nil {
		if err = userRepo.Save(&userInfo); err != nil {
			return util.ResponseError(c, http.StatusInternalServerError, "500-100", "failed to save user", err)
		}

		if user, err = userRepo.FindByEmail(userInfo.Email); err != nil {
			return util.ResponseError(c, http.StatusInternalServerError, "500-109", "failed to get user", err)
		}
	}

	return c.JSON(http.StatusOK, user)
}
