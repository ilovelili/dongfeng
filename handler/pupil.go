package handler

import (
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetPupils GET /pupils
func GetPupils(c echo.Context) error {
	class, year := c.QueryParam("class"), c.QueryParam("year")
	pupils, err := pupilRepo.Find(class, year)
	if err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-107", "failed to get pupils", err)
	}

	return c.JSON(http.StatusOK, pupils)
}

// UpdatePupils POST /pupils
func UpdatePupils(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-105", "failed to parse pupils", err)
	}
	defer file.Close()

	pupils := []*model.Pupil{}
	if err := gocsv.Unmarshal(file, &pupils); err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-105", "failed to parse pupils", err)
	}

	for _, pupil := range pupils {
		pupil.CreatedBy = userInfo.Email
	}

	if err := pupilRepo.DeleteInsert(pupils); err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-108", "failed to save pupils", err)
	}

	notify(model.PupilListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
