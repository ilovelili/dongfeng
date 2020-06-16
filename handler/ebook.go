package handler

import (
	"net/http"
	"strconv"

	"github.com/ilovelili/dongfeng/core/controller"
	"github.com/ilovelili/dongfeng/core/model"
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

// UpdateEbook POST /ebook
func UpdateEbook(c echo.Context) error {
	ebook := new(model.Ebook)
	if err := c.Bind(ebook); err != nil {
		return util.ResponseError(c, "400-119", "failed to bind ebook", err)
	}

	ebookCtrl := controller.NewEbookController()
	if err := ebookCtrl.SaveEbook(ebook); err != nil {
		return util.ResponseError(c, "500-131", "failed to generate ebook", err)
	}

	return c.NoContent(http.StatusOK)
}

// UpdateProfileTemplatePreview POST /convertProfileToTemplate
func UpdateProfileTemplatePreview(c echo.Context) error {
	templatePreview := new(model.TemplatePreview)
	if err := c.Bind(templatePreview); err != nil {
		return util.ResponseError(c, "400-120", "failed to bind template preview", err)
	}

	ebookCtrl := controller.NewEbookController()
	if err := ebookCtrl.SaveTemplatePreview(templatePreview); err != nil {
		return util.ResponseError(c, "500-131", "failed to generate ebook", err)
	}

	return c.NoContent(http.StatusOK)
}
