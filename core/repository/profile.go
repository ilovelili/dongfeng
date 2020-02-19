package repository

import (
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/jinzhu/gorm"
)

// Profile repository
type Profile struct{}

// NewProfileRepository profile repository
func NewProfileRepository() *Profile {
	db().AutoMigrate(&model.ProfileTemplate{}, &model.Profile{})
	return new(Profile)
}

// FindTemplates find all templates
func (r *Profile) FindTemplates() ([]*model.ProfileTemplate, error) {
	templates := []*model.ProfileTemplate{}
	err := db().Find(&templates).Error
	return templates, err
}

// FindTemplateByName find template by name
func (r *Profile) FindTemplateByName(name string) (*model.ProfileTemplate, error) {
	template := new(model.ProfileTemplate)
	err := db().Where("name = ?", name).First(&template).Error
	return template, err
}

// SaveTemplate save template
func (r *Profile) SaveTemplate(template *model.ProfileTemplate) error {
	_template := new(model.ProfileTemplate)
	if err := db().Where("name = ?", template.Name).Find(&_template).Error; err == nil {
		template.ID = _template.ID
	}

	return db().Save(template).Error
}

// DeleteTemplate delete template
func (r *Profile) DeleteTemplate(name string) error {
	template := new(model.ProfileTemplate)
	if err := db().Where("name = ?", name).First(&template).Error; err != nil {
		return err
	}

	return db().Delete(template).Error
}

// FindProfiles find profiles
func (r *Profile) FindProfiles(year string) ([]*model.Profile, error) {
	profiles := []*model.Profile{}
	err := db().
		Joins("JOIN pupils ON profiles.pupil_id = pupils.id").Joins("JOIN classes ON pupils.class_id = classes.id").Where("classes.year = ?", year).
		Preload("Template").Preload("Pupil").Preload("Pupil.Class").
		Find(&profiles).Error
	return profiles, err
}

// FindProfile find profile by id
func (r *Profile) FindProfile(id string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("id = ?", id).
		Preload("Template").Preload("Pupil").Preload("Pupil.Class").
		First(&profile).Error
	return profile, err
}

// FindPrevProfile find previous profile
func (r *Profile) FindPrevProfile(pupilID, date string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("pupil_id = ? AND date < ?", pupilID, date).
		Preload("Template").Preload("Pupil").Preload("Pupil.Class").
		Order("date desc").
		First(&profile).Error

	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return profile, err
}

// FindNextProfile find next profile
func (r *Profile) FindNextProfile(pupilID, date string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("pupil_id = ? AND date > ?", pupilID, date).
		Preload("Template").Preload("Pupil").Preload("Pupil.Class").
		Order("date desc").
		First(&profile).Error

	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return profile, err
}

// SaveProfile save profile
func (r *Profile) SaveProfile(profile *model.Profile) error {
	if profile.ID == 0 {
		_profile := new(model.Profile)
		if err := db().Where("pupil_id AND date = ?", profile.PupilID, profile.Date).Find(&_profile).Error; err == nil {
			profile.ID = _profile.ID
		}
		return db().Save(profile).Error
	}
	return db().Model(&model.Profile{}).Update("profile", profile.Profile).Error
}

// DeleteProfile delete profile
func (r *Profile) DeleteProfile(id string) error {
	return db().Unscoped().Where("id = ?", id).Delete(&model.Profile{}).Error
}
