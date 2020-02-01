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
		return util.ResponseError(c, "500-107", "failed to get pupils", err)
	}

	return c.JSON(http.StatusOK, pupils)
}

// SavePupils POST /pupils
func SavePupils(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, "400-105", "failed to parse pupils", err)
	}
	defer file.Close()

	pupils := []*model.Pupil{}
	if err := gocsv.Unmarshal(file, &pupils); err != nil {
		return util.ResponseError(c, "400-105", "failed to parse pupils", err)
	}

	// year|classname|classid map
	classMap := make(map[string]map[string]uint)
	for _, pupil := range pupils {
		year := pupil.CSVYear
		if _, ok := classMap[year]; !ok {
			classes, err := classRepo.Find(year)
			if err != nil {
				return util.ResponseError(c, "500-105", "failed to get classes", err)
			}

			classMap[year] = make(map[string]uint)
			for _, class := range classes {
				classMap[year][class.Name] = class.ID
			}
		}

	CLASSLOOP:
		for k, v := range classMap {
			if year == k {
				for className, classID := range v {
					if className == pupil.CSVClass {
						pupil.ClassID = classID
						break CLASSLOOP
					}
				}
			}
		}

		pupil.CreatedBy = userInfo.Email
	}

	if err := pupilRepo.DeleteInsert(pupils); err != nil {
		return util.ResponseError(c, "500-108", "failed to save pupils", err)
	}

	notify(model.PupilListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
