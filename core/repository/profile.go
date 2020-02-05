package repository

import "github.com/ilovelili/dongfeng/core/model"

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
	err := db().First(&template).Error
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

// FindAllProfiles find all profiles
func (r *Profile) FindAllProfiles() ([]*model.Profile, error) {
	profiles := []*model.Profile{}
	err := db().Find(&profiles).Error
	return profiles, err
}

// FindProfile find profile
func (r *Profile) FindProfile(pupilID uint, date string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("pupil_id = ? AND date = ?", pupilID, date).First(&profile).Error
	return profile, err
}

// FindPrevProfile find previous profile
func (r *Profile) FindPrevProfile(pupilID uint, date string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("pupil_id = ? AND date < ? ORDER BY date desc", pupilID, date).First(&profile).Error
	return profile, err
}

// FindNextProfile find next profile
func (r *Profile) FindNextProfile(pupilID uint, date string) (*model.Profile, error) {
	profile := new(model.Profile)
	err := db().Where("pupil_id = ? AND date > ? ORDER BY date desc", pupilID, date).First(&profile).Error
	return profile, err
}

// SaveProfiles save profiles
func (r *Profile) SaveProfiles(profiles []*model.Profile) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, profile := range profiles {
		_profile := new(model.Profile)
		if err := db().Where("pupil_id = ? AND date = ?", profile.PupilID, profile.Date).Find(&_profile).Error; err == nil {
			profile.ID = _profile.ID
		}

		if err := tx.Save(profile).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteProfiles delete profiles
func (r *Profile) DeleteProfiles(profiles []*model.Profile) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, profile := range profiles {
		_profile := new(model.Profile)
		if err := db().Where("pupil_id = ? AND date = ?", profile.PupilID, profile.Date).First(&_profile).Error; err != nil {
			tx.Rollback()
			return err
		}

		profile.ID = _profile.ID
		if err := tx.Delete(profile).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
