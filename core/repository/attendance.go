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

// Delete soft delete absence
func (r *Attendance) Delete(absence *model.Absence) error {
	return db().Model(&model.Absence{}).Delete(absence).Error
}
