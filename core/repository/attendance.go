package repository

import (
	"github.com/ilovelili/dongfeng/core/model"
)

// Attendance attendance repository
type Attendance struct{}

// NewAttendanceRepository new attendance repository
func NewAttendanceRepository() *Attendance {
	db().AutoMigrate(&model.Absence{}, &model.Holiday{}, &model.Pupil{})
	return new(Attendance)
}

// FindHolidays find holidays
func (r *Attendance) FindHolidays() ([]*model.Holiday, error) {
	holidays := []*model.Holiday{}
	err := db().Find(&holidays).Error
	return holidays, err
}

// FindAbsences find absences
func (r *Attendance) FindAbsences(from, to string) ([]*model.Absence, error) {
	absences := []*model.Absence{}
	query := db().Joins("JOIN pupils ON absences.pupil_id = pupils.id")
	if from != "" && to != "" {
		query = query.Where("absences.date BETWEEN ? AND ?", from, to)
	} else if from != "" && to == "" {
		query = query.Where("absences.date > ?", from)
	} else if from == "" && to != "" {
		query = query.Where("absences.date < ?", to)
	}

	err := query.Preload("Pupil").Find(&absences).Error
	return absences, err
}

// FindAbsencesByPupil find absences by pupil
func (r *Attendance) FindAbsencesByPupil(pupilID uint, from, to string) ([]*model.Absence, error) {
	absences := []*model.Absence{}
	query := db().Joins("JOIN pupils ON absences.pupil_id = pupils.id AND pupils.id = ?", pupilID)
	if from != "" && to != "" {
		query = query.Where("absences.date BETWEEN ? AND ?", from, to)
	} else if from != "" && to == "" {
		query = query.Where("absences.date > ?", from)
	} else if from == "" && to != "" {
		query = query.Where("absences.date < ?", to)
	}

	err := query.Preload("Pupil").Find(&absences).Error
	return absences, err
}

// FindAbsencesByClass find absences by class
func (r *Attendance) FindAbsencesByClass(classID uint, from, to string) ([]*model.Absence, error) {
	absences := []*model.Absence{}
	query := db().Joins("JOIN pupils ON absences.pupil_id = pupils.id AND pupils.class_id = ?", classID)
	if from != "" && to != "" {
		query = query.Where("absences.date BETWEEN ? AND ?", from, to)
	} else if from != "" && to == "" {
		query = query.Where("absences.date > ?", from)
	} else if from == "" && to != "" {
		query = query.Where("absences.date < ?", to)
	}

	err := query.Preload("Pupil").Find(&absences).Error
	return absences, err
}

// Save save absence
func (r *Attendance) Save(absence *model.Absence) error {
	return db().Model(&model.Absence{}).Save(absence).Error
}

// SaveAll save all absences, skip the items for already saved
func (r *Attendance) SaveAll(absences []*model.Absence) error {
	candidates := []*model.Absence{}
	for _, absence := range absences {
		// if already there, skip
		_absence := new(model.Absence)
		if db().Where("absences.date = ? AND absences.pupil_id = ?", absence.Date, absence.PupilID).Find(&_absence).RecordNotFound() {
			candidates = append(candidates, absence)
		}
	}

	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, candidate := range candidates {
		// set ID to 0 to insert instead of update
		candidate.ID = 0
		if err := tx.Save(candidate).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Delete hard delete absence
func (r *Attendance) Delete(absence *model.Absence) error {
	return db().Unscoped().Model(&model.Absence{}).Delete(absence).Error
}
