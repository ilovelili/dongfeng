package handler

import (
	"net/http"
	"strconv"

	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetEbooks GET /ebooks
func GetEbooks(c echo.Context) error {
	var (
		err     error
		ebooks  interface{}
		pupilID uint = 0
		classID uint = 0
		year         = c.QueryParam("year")
		class        = c.QueryParam("class")
		pupil        = c.QueryParam("name")
	)

	if class != "" {
		_classID, err := strconv.ParseUint(class, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-108", "invalid class", err)
		}
		classID = uint(_classID)
	}
	if pupil != "" {
		_pupilID, err := strconv.ParseUint(pupil, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-109", "invalid pupil", err)
		}
		pupilID = uint(_pupilID)
	}

	if pupilID != 0 {
		ebooks, err = ebookRepo.FindByPupilID(pupilID)
	} else if classID != 0 {
		ebooks, err = ebookRepo.FindByClassID(classID)
	} else {
		ebooks, err = ebookRepo.Find(year)
	}

	if err != nil {
		return util.ResponseError(c, "500-124", "failed to get ebooks", err)
	}
	return c.JSON(http.StatusOK, ebooks)
}
