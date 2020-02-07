package handler

import (
	"io/ioutil"
	"net/http"
	"strconv"

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
	if template.Profile == nil {
		return c.JSONBlob(http.StatusOK, []byte{})
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

// GetProfiles GET /profiles
func GetProfiles(c echo.Context) error {
	year := c.QueryParam("year")
	profiles, err := profileRepo.FindProfiles(year)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	// omit profile content
	_profiles := []*model.Profile{}
	for _, profile := range profiles {
		_profiles = append(_profiles, &model.Profile{
			BaseModel:  profile.BaseModel,
			Pupil:      profile.Pupil,
			PupilID:    profile.PupilID,
			Template:   profile.Template,
			TemplateID: profile.TemplateID,
			Date:       profile.Date,
			CreatedBy:  profile.CreatedBy,
		})
	}
	return c.JSON(http.StatusOK, _profiles)
}

// GetPreviousProfile GET /profile/prev
func GetPreviousProfile(c echo.Context) error {
	pupilID, date := c.QueryParam("name"), c.QueryParam("date")
	profile, err := profileRepo.FindPrevProfile(pupilID, date)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	// omit profile content
	return c.JSON(http.StatusOK, &model.Profile{
		BaseModel:  profile.BaseModel,
		Pupil:      profile.Pupil,
		PupilID:    profile.PupilID,
		Template:   profile.Template,
		TemplateID: profile.TemplateID,
		Date:       profile.Date,
		CreatedBy:  profile.CreatedBy,
	})
}

// GetNextProfile GET /profile/next
func GetNextProfile(c echo.Context) error {
	pupilID, date := c.QueryParam("name"), c.QueryParam("date")
	profile, err := profileRepo.FindNextProfile(pupilID, date)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	// omit profile content
	return c.JSON(http.StatusOK, &model.Profile{
		BaseModel:  profile.BaseModel,
		Pupil:      profile.Pupil,
		PupilID:    profile.PupilID,
		Template:   profile.Template,
		TemplateID: profile.TemplateID,
		Date:       profile.Date,
		CreatedBy:  profile.CreatedBy,
	})
}

// SaveProfile POST /profile
func SaveProfile(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	profile := new(model.Profile)
	if err := c.Bind(profile); err != nil {
		return util.ResponseError(c, "400-118", "failed to bind profile", err)
	}

	profile.CreatedBy = userInfo.Email
	err := profileRepo.SaveProfile(profile)
	if err != nil {
		return util.ResponseError(c, "500-129", "failed to save profiles", err)
	}

	notify(model.GrowthProfileUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// DeleteProfile DELETE /profile
func DeleteProfile(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	id := c.QueryParam("id")
	err := profileRepo.DeleteProfile(id)
	if err != nil {
		return util.ResponseError(c, "500-130", "failed to delete profiles", err)
	}

	notify(model.GrowthProfileUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// GetProfileContent GET /profileContent
func GetProfileContent(c echo.Context) error {
	id := c.QueryParam("id")
	profile, err := profileRepo.FindProfile(id)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}
	if profile.Profile == nil {
		if profile.Template == nil {
			return c.JSONBlob(http.StatusOK, []byte{})
		}
		return c.JSONBlob(http.StatusOK, []byte(*profile.Template.Profile))
	}
	return c.JSONBlob(http.StatusOK, []byte(*profile.Profile))
}

// SaveProfileContent POST /profileContent
func SaveProfileContent(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	id := c.QueryParam("id")
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return util.ResponseError(c, "400-108", "invalid class", err)
	}

	profile := new(model.Profile)
	profile.ID = uint(_id)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return util.ResponseError(c, "400-116", "invalid profile template", err)
	}
	content := string(body)
	profile.Profile = &content
	profile.CreatedBy = userInfo.Email

	if err := profileRepo.SaveProfile(profile); err != nil {
		return util.ResponseError(c, "500-129", "failed to save profiles", err)
	}

	return c.NoContent(http.StatusOK)
}
