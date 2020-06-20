package model

// ProfileTemplate  profile template
type ProfileTemplate struct {
	BaseModel
	Name      string  `json:"name"`
	Profile   *string `sql:"size:999999" json:"profile"`
	Tags      *string `json:"tags"`
	CreatedBy string  `json:"created_by"`
}

// Profile profile entity
type Profile struct {
	BaseModel
	Pupil      *Pupil           `json:"pupil"`
	PupilID    *uint            `json:"pupil_id"`
	Template   *ProfileTemplate `json:"profile_template"`
	TemplateID uint             `json:"template_id"`
	Date       string           `json:"date"`
	Profile    *string          `sql:"size:999999" json:"profile"`
	CreatedBy  string           `json:"created_by"`
}

// ProfileCount profile count entity just for json response
type ProfileCount struct {
	PupilID uint   `gorm:"column:pupilId" json:"pupilId"`
	Pupil   string `json:"pupil"`
	Class   string `json:"class"`
}
