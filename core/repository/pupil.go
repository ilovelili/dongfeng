package repository

import "github.com/ilovelili/dongfeng/core/model"

// Pupil pupil repository
type Pupil struct{}

// NewPupilRepository new pupil repository
func NewPupilRepository() *Pupil {
	db().AutoMigrate(&model.Pupil{})
	return new(Pupil)
}

// FindAll find all pupiles
func (r *Pupil) FindAll() ([]*model.Pupil, error) {
	pupiles := []*model.Pupil{}
	err := db().Preload("Class").Find(&pupiles).Error
	return pupiles, err
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
		if err := tx.Save(pupil).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
