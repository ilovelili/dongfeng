package repository

import "github.com/ilovelili/dongfeng/core/model"

// Class class repository
type Class struct{}

// NewClassRepository new class repository
func NewClassRepository() *Class {
	db().AutoMigrate(&model.Class{})
	return new(Class)
}

// FindAll find all classes
func (r *Class) FindAll(year string) ([]*model.Class, error) {
	classes := []*model.Class{}
	query := db()
	if year != "" {
		query = query.Where("year = ?", year)
	}
	err := query.Find(&classes).Error
	return classes, err
}

// DeleteInsert delete and insert classes
func (r *Class) DeleteInsert(classes []*model.Class) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	years := []string{}
	for _, class := range classes {
		yearFound := false
		for _, year := range years {
			if year == class.Year {
				yearFound = true
			}
		}

		if !yearFound {
			years = append(years, class.Year)
		}
	}

	// use unscoped for physical delete
	if err := tx.Unscoped().Where("year IN (?)", years).Delete(&model.Class{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, class := range classes {
		if err := tx.Save(class).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
