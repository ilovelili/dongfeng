package repository

import "github.com/ilovelili/dongfeng/core/model"

// Teacher teacher repository
type Teacher struct{}

// NewTeacherRepository new teacher repository
func NewTeacherRepository() *Teacher {
	db().AutoMigrate(&model.Teacher{}, &model.User{}, &model.User{})
	return new(Teacher)
}

// Find find teachers
func (r *Teacher) Find(class, year string) ([]*model.Teacher, error) {
	teachers := []*model.Teacher{}
	query := db().Joins("LEFT JOIN classes ON classes.id = teachers.class_id").Joins("LEFT JOIN users ON users.id = teachers.user_id")
	if class != "" && year != "" {
		query = query.Where("teachers.class_id = 0 OR (classes.year = ? AND classes.name = ?)", year, class)
	} else if class == "" && year != "" {
		query = query.Where("teachers.class_id = 0 OR classes.year = ?", year)
	} else if class != "" && year == "" {
		query = query.Where("teachers.class_id = 0 OR classes.name = ?", class)
	}

	err := query.Preload("Class").Preload("User").Find(&teachers).Error
	return teachers, err
}

// Save save teacher
func (r *Teacher) Save(teacher *model.Teacher) error {
	return db().Model(&model.Teacher{}).Where("id = ?", teacher.ID).Update("class_id", teacher.ClassID, "name", teacher.Name).Error
}

// DeleteInsert delete and insert teachers
func (r *Teacher) DeleteInsert(teachers []*model.Teacher) error {
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
	for _, teacher := range teachers {
		classIDFound := false
		for _, classID := range classIDs {
			if teacher.ClassID != nil && *teacher.ClassID == classID {
				classIDFound = true
			}
		}

		if !classIDFound {
			classIDs = append(classIDs, *teacher.ClassID)
		}
	}

	// use unscoped for physical delete
	if err := tx.Unscoped().Where("class_id IN (?)", classIDs).Delete(&model.Teacher{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, teacher := range teachers {
		// set ID to 0 to insert instead of update
		teacher.ID = 0
		if err := tx.Save(teacher).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
