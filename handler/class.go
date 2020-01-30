package handler

import (
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetClasses GET /classes
func GetClasses(c echo.Context) error {
	year := c.QueryParam("year")
	classes, err := classRepo.Find(year)
	if err != nil {
		return util.ResponseError(c, "500-105", "failed to get classes", err)
	}

	return c.JSON(http.StatusOK, classes)
}

// SaveClasses POST /classes
func SaveClasses(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, "400-104", "failed to parse classes", err)
	}
	defer file.Close()

	classes := []*model.Class{}
	if err := gocsv.Unmarshal(file, &classes); err != nil {
		return util.ResponseError(c, "400-104", "failed to parse classes", err)
	}

	for _, class := range classes {
		class.CreatedBy = userInfo.Email
	}

	if err := classRepo.DeleteInsert(classes); err != nil {
		return util.ResponseError(c, "500-106", "failed to save classes", err)
	}

	notify(model.ClassListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
