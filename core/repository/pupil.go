package repository

import "github.com/ilovelili/dongfeng/core/model"

// Pupil pupil repository
type Pupil struct{}

// NewPupilRepository new pupil repository
func NewPupilRepository() *Pupil {
	db().AutoMigrate(&model.Pupil{}, &model.Class{})
	return new(Pupil)
}

// Find find pupils
func (r *Pupil) Find(class, year string) ([]*model.Pupil, error) {
	pupils := []*model.Pupil{}

	query := db().Joins("JOIN classes ON classes.id = pupils.class_id")
	if class != "" && year != "" {
		query = query.Where("classes.year = ? AND classes.id = ?", year, class)
	} else if class == "" && year != "" {
		query = query.Where("classes.year = ?", year)
	} else if class != "" && year == "" {
		query = query.Where("classes.id = ?", class)
	}

	err := query.Preload("Class").Find(&pupils).Error
	return pupils, err
}

// FindByPupilID find pupils by pupilID
func (r *Pupil) FindByPupilID(pupilID uint) ([]*model.Pupil, error) {
	pupils := []*model.Pupil{}
	err := db().Where("pupils.id = ?", pupilID).Preload("Class").Find(&pupils).Error
	return pupils, err
}

// FindByClassID find pupils by classID
func (r *Pupil) FindByClassID(classID uint) ([]*model.Pupil, error) {
	pupils := []*model.Pupil{}
	err := db().Where("pupils.class_id = ?", classID).Preload("Class").Find(&pupils).Error
	return pupils, err
}

// DeleteInsert delete and insert pupils
func (r *Pupil) DeleteInsert(pupils []*model.Pupil) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	classIDs := []uint{}
	for _, pupil := range pupils {
		classIDFound := false
		for _, classID := range classIDs {
			if classID == pupil.ClassID {
				classIDFound = true
			}
		}

		if !classIDFound {
			classIDs = append(classIDs, pupil.ClassID)
		}
	}

	if err := tx.Unscoped().Where("class_id IN (?)", classIDs).Delete(&model.Pupil{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, pupil := range pupils {
		// set ID to 0 to insert instead of update
		pupil.ID = 0
		if err := tx.Save(pupil).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
