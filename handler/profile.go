package handler

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ilovelili/dongfeng/core/controller"
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

// SaveProfileTemplateTags POST /profileTemplate/tags
func SaveProfileTemplateTags(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	name := c.QueryParam("name")
	tags := c.QueryParam("tags")

	template, err := profileRepo.FindTemplateByName(name)
	if err != nil {
		template = &model.ProfileTemplate{Name: name, Tags: &tags}
	} else {
		template.Tags = &tags
	}

	template.CreatedBy = userInfo.Email
	if err := profileRepo.UpdateTemplateTags(template); err != nil {
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

// GetProfileTemplateTags GET /profileTemplateTags
func GetProfileTemplateTags(c echo.Context) error {
	name := c.QueryParam("name")
	template, err := profileRepo.FindTemplateByName(name)
	if err != nil {
		return util.ResponseError(c, "500-125", "failed to get profile templates", err)
	}
	return c.JSON(http.StatusOK, template.Tags)
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

	profile, err = profileRepo.FindProfile(strconv.Itoa(int(profile.ID)))
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	return c.JSON(http.StatusOK, profile)
}

// DeleteProfile DELETE /profile
func DeleteProfile(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	id := c.QueryParam("id")

	profile, err := profileRepo.FindProfile(id)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	err = profileRepo.DeleteProfile(id)
	if err != nil {
		return util.ResponseError(c, "500-130", "failed to delete profiles", err)
	}

	// try to delete original ebook file if exists
	ebook := &model.Ebook{
		Pupil: profile.Pupil,
		Date:  profile.Date,
	}

	ebookCtrl := controller.NewEbookController()
	if err := ebookCtrl.RemoveFromStorage(ebook); err != nil {
		return util.ResponseError(c, "500-133", "failed to delete ebook", err)
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

// GetProfileCount GET /profileCount
func GetProfileCount(c echo.Context) error {
	classID := c.QueryParam("classId")
	date := c.QueryParam("date")
	count, err := profileRepo.GetProfileByClassIDAndDate(classID, date)
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}
	return c.JSON(http.StatusOK, count)
}

// SaveProfileContent POST /profileContent
func SaveProfileContent(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	id := c.QueryParam("id")
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return util.ResponseError(c, "400-108", "invalid id", err)
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

// ConvertProfileToTemplate POST /convertProfileToTemplate
func ConvertProfileToTemplate(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	id := c.QueryParam("id")
	templateName := c.QueryParam("templateName")
	tags := c.QueryParam("tags")

	_, err := profileRepo.FindTemplateByName(templateName)
	// tamplate found ... template name already exsits
	if err == nil {
		return util.ResponseError(c, "500-132", "profile template already exists", err)
	}

	profile, err := profileRepo.FindProfile(id)
	// profile not found
	if err != nil {
		return util.ResponseError(c, "500-128", "failed to get profiles", err)
	}

	if err := profileRepo.SaveTemplate(&model.ProfileTemplate{
		Name:      templateName,
		CreatedBy: userInfo.Email,
		Profile:   profile.Profile,
		Tags:      &tags,
	}); err != nil {
		return util.ResponseError(c, "500-126", "failed to save profile templates", err)
	}

	return c.NoContent(http.StatusOK)
}

// ApplyProfileTemplate POST /applyProfileTemplate
func ApplyProfileTemplate(c echo.Context) error {
	type ApplyProfileRequest struct {
		Date       string `json:"date"`
		ClassID    uint   `json:"classId"`
		PupilIDs   []uint `json:"pupilIds"`
		TemplateID uint   `json:"templateId"`
		Overwrite  bool   `json:"overwrite"`
	}

	userInfo, _ := c.Get("userInfo").(model.User)
	applyProfileReq := new(ApplyProfileRequest)
	if err := c.Bind(applyProfileReq); err != nil {
		return util.ResponseError(c, "400-121", "failed to bind apply profile template request", err)
	}

	// first, get pupils in class
	pupils, err := pupilRepo.FindByClassID(applyProfileReq.ClassID)
	if err != nil {
		return util.ResponseError(c, "500-107", "failed to get pupils", err)
	}

	// then filter out the pupil ids
	pupilIDs := []uint{}
	for _, pupil := range pupils {
		// if overwrite set to false, then exclude the pupilIds
		if !applyProfileReq.Overwrite {
			for _, pupilID := range applyProfileReq.PupilIDs {
				if pupilID == pupil.ID {
					continue
				}
			}
		}
		pupilIDs = append(pupilIDs, pupil.ID)
	}

	// apply
	err = profileRepo.ApplyProfileTemplate(applyProfileReq.Date, pupilIDs, applyProfileReq.TemplateID, userInfo.Email)
	if err != nil {
		return util.ResponseError(c, "500-135", "failed to apply template", err)
	}

	return c.NoContent(http.StatusOK)
}
