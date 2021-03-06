package repository

import "github.com/ilovelili/dongfeng/core/model"

// Ebook ebook repository
type Ebook struct{}

// NewEbookRepository new ebook repository
func NewEbookRepository() *Ebook {
	db().AutoMigrate(&model.Ebook{}, &model.TemplatePreview{})
	return new(Ebook)
}

// Find find ebooks
func (r *Ebook) Find(year string) ([]*model.Ebook, error) {
	ebooks := []*model.Ebook{}
	query := db().Joins("JOIN pupils ON ebooks.pupil_id = pupils.id").Joins("JOIN classes ON pupils.class_id = classes.id")
	if year != "" {
		query = query.Where("classes.year = ? AND ebooks.converted = 1", year)
	} else {
		query = query.Where("ebooks.converted = 1")
	}

	err := query.Preload("Pupil").Preload("Pupil.Class").Find(&ebooks).Error
	return ebooks, err
}

// FindByID find by id
func (r *Ebook) FindByID(id uint) (*model.Ebook, error) {
	ebook := new(model.Ebook)
	err := db().Where("id = ?", id).Preload("Pupil").Preload("Pupil.Class").Find(&ebook).Error
	return ebook, err
}

// FindByPupilID find by pupil id
func (r *Ebook) FindByPupilID(pupilID uint) ([]*model.Ebook, error) {
	ebooks := []*model.Ebook{}
	err := db().Where("converted = 1 AND pupil_id = ?", pupilID).Preload("Pupil").Preload("Pupil.Class").Find(&ebooks).Error
	return ebooks, err
}

// FindByClassID find by class id
func (r *Ebook) FindByClassID(classID uint) ([]*model.Ebook, error) {
	ebooks := []*model.Ebook{}
	err := db().Joins("JOIN pupils ON ebooks.pupil_id = pupils.id").Where("ebooks.converted = 1 AND pupils.class_id = ?", classID).Preload("Pupil").Preload("Pupil.Class").Find(&ebooks).Error
	return ebooks, err
}

// Save save ebook
func (r *Ebook) Save(ebook *model.Ebook, forceUpdate bool) (dirty bool, err error) {
	_ebook := new(model.Ebook)
	err = db().Where("pupil_id = ? AND date = ?", ebook.PupilID, ebook.Date).First(&_ebook).Error
	if err != nil {
		dirty = true
		err = db().Save(ebook).Error
	} else if _ebook.Hash == ebook.Hash {
		dirty = false
		ebook.ID = _ebook.ID
		// force update
		if forceUpdate {
			err = db().Save(ebook).Error
		}
	} else {
		dirty = true
		ebook.ID = _ebook.ID
		err = db().Save(ebook).Error
	}

	return
}

// SaveTemplatePreview save template preview
func (r *Ebook) SaveTemplatePreview(templatePreview *model.TemplatePreview, forceUpdate bool) (dirty bool, err error) {
	_templatePreview := new(model.TemplatePreview)
	err = db().Where("name = ?", templatePreview.Name).First(&_templatePreview).Error
	if err != nil {
		dirty = true
		err = db().Save(templatePreview).Error
	} else if _templatePreview.Hash == templatePreview.Hash {
		dirty = false
		templatePreview.ID = _templatePreview.ID
		// force update
		if forceUpdate {
			err = db().Save(templatePreview).Error
		}
	} else {
		dirty = true
		templatePreview.ID = _templatePreview.ID
		err = db().Save(templatePreview).Error
	}

	return
}
