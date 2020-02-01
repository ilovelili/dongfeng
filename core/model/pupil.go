package model

// Pupil pupil entity
type Pupil struct {
	BaseModel
	Name      string `json:"name" csv:"姓名"`
	CSVYear   string `gorm:"-" json:"-" csv:"学年"`
	CSVClass  string `gorm:"-" json:"-" csv:"班级"`
	Class     Class  `json:"class" csv:"-"`
	ClassID   uint   `json:"class_id" csv:"-"`
	CreatedBy string `json:"created_by" csv:"-"`
}
