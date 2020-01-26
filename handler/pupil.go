package handler

import (
	"net/http"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetPupils GET /pupils
func GetPupils(c echo.Context) error {
	pupils, err := pupilRepo.FindAll()
	if err != nil {

		return util.ResponseError(c, http.StatusInternalServerError, "500-107", "failed to get pupiles", err)
	}

	return c.JSON(http.StatusOK, pupils)
}

// UpdatePupils POST /pupils
func UpdatePupils(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	pupils := []*model.Pupil{}
	if err := c.Bind(pupils); err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-105", "failed to bind pupils", err)
	}

	if err := pupilRepo.DeleteInsert(pupils); err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-108", "failed to save pupils", err)
	}

	notify(model.PupilListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
