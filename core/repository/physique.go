package repository

import (
	"github.com/ilovelili/dongfeng/core/model"
)

// Physique pupil repository
type Physique struct{}

// NewPhysiqueRepository new pupil repository
func NewPhysiqueRepository() *Physique {
	db().AutoMigrate(
		&model.Physique{},
		&model.AgeHeightWeightPMaster{},
		&model.AgeHeightWeightSDMaster{},
		&model.BMIMaster{},
		&model.HeightToWeightPMaster{},
		&model.HeightToWeightSDMaster{},
		&model.Pupil{},
	)
	return new(Physique)
}

// SelectAgeHeightWeightPMasters select age weight / height p master
func (r *Physique) SelectAgeHeightWeightPMasters() ([]*model.AgeHeightWeightPMaster, error) {
	masters := []*model.AgeHeightWeightPMaster{}
	err := db().Find(&masters).Error
	return masters, err
}

// SelectAgeHeightWeightSDMasters select age weight / height sd master
func (r *Physique) SelectAgeHeightWeightSDMasters() ([]*model.AgeHeightWeightSDMaster, error) {
	masters := []*model.AgeHeightWeightSDMaster{}
	err := db().Find(&masters).Error
	return masters, err
}

// SelectHeightToWeightPMasters select height weight p masters
func (r *Physique) SelectHeightToWeightPMasters() ([]*model.HeightToWeightPMaster, error) {
	masters := []*model.HeightToWeightPMaster{}
	err := db().Find(&masters).Error
	return masters, err
}

// SelectHeightToWeightSDMasters select height to weight sd masters
func (r *Physique) SelectHeightToWeightSDMasters() ([]*model.HeightToWeightSDMaster, error) {
	masters := []*model.HeightToWeightSDMaster{}
	err := db().Find(&masters).Error
	return masters, err
}

// SelectBMIMasters select bmi sd master
func (r *Physique) SelectBMIMasters() ([]*model.BMIMaster, error) {
	masters := []*model.BMIMaster{}
	err := db().Find(&masters).Error
	return masters, err
}

// Find physique by pupil ID
func (r *Physique) Find(pupilIDs []uint) ([]*model.Physique, error) {
	physiques := []*model.Physique{}
	err := db().Where("physiques.pupil_id IN (?)", pupilIDs).Preload("Pupil").Preload("Pupil.Class").Find(&physiques).Error
	return physiques, err
}

// FindAll find all physiques by year
func (r *Physique) FindAll(year string) ([]*model.Physique, error) {
	physiques := []*model.Physique{}
	err := db().
		Joins("JOIN pupils ON physiques.pupil_id = pupils.id").
		Joins("JOIN classes ON pupils.class_id = classes.id AND classes.year = ?", year).
		Preload("Pupil").Preload("Pupil.Class").
		Find(&physiques).Error

	return physiques, err
}

// Save save physique
func (r *Physique) Save(physique *model.Physique) error {
	return db().Save(physique).Error
}

// SaveAll save all physiques
func (r *Physique) SaveAll(physiques []*model.Physique) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, physique := range physiques {
		_physique := new(model.Physique)
		if tx.Where("physiques.pupil_id = ?", physique.ID).Find(&_physique).RecordNotFound() {
			// if not found, insert. else update
			physique.ID = 0
		}

		if err := tx.Save(physique).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
