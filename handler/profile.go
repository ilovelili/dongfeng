package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetProfileTemplate GET /profileTemplate
func GetProfileTemplate(c echo.Context) error {
	name := c.QueryParam("name")
	template, err := profileRepo.FindTemplateByName(name)
	if err != nil {
		return util.ResponseError(c, "500-125", "failed to get profile templates", err)
	}
	return c.JSON(http.StatusOK, &model.ProfileTemplate{
		BaseModel: template.BaseModel,
		Name:      template.Name,
		CreatedBy: template.CreatedBy,
	})
}

// GetProfileTemplateContent GET /profileTemplateContent
func GetProfileTemplateContent(c echo.Context) error {
	name := c.QueryParam("name")
	template, err := profileRepo.FindTemplateByName(name)
	if err != nil {
		return util.ResponseError(c, "500-125", "failed to get profile templates", err)
	}
	return c.JSONBlob(http.StatusOK, []byte(*template.Profile))
}

// SaveProfileTemplate POST /profileTemplate
func SaveProfileTemplate(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	name := c.QueryParam("name")
	template, err := profileRepo.FindTemplateByName(name)
	if err != nil {
		template = &model.ProfileTemplate{Name: name}
	} else {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return util.ResponseError(c, "400-116", "invalid profile template", err)
		}
		profile := string(body)
		template.Profile = &profile
	}
	template.CreatedBy = userInfo.Email
	if err := profileRepo.SaveTemplate(template); err != nil {
		return util.ResponseError(c, "500-126", "failed to save profile templates", err)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteProfileTemplate DELETE /profileTemplate
func DeleteProfileTemplate(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	name := c.QueryParam("name")

	if err := profileRepo.DeleteTemplate(name); err != nil {
		return util.ResponseError(c, "500-127", "failed to delete profile templates", err)
	}

	notify(model.GrowthProfileTemplateUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// GetProfileTemplates GET /profileTemplates
func GetProfileTemplates(c echo.Context) error {
	templates, err := profileRepo.FindTemplates()
	if err != nil {
		return util.ResponseError(c, "500-125", "failed to get profile templates", err)
	}
	return c.JSON(http.StatusOK, templates)
}
