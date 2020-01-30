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
func (r *Physique) Find(pupil string) ([]*model.Physique, error) {
	physiques := []*model.Physique{}
	err := db().Where("physiques.pupil_id = ?", pupil).Find(&physiques).Error
	return physiques, err
}

// FindAll find all physiques
func (r *Physique) FindAll() ([]*model.Physique, error) {
	physiques := []*model.Physique{}
	err := db().Find(&physiques).Error
	return physiques, err
}
