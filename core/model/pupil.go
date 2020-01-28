package model

// Pupil pupil entity
type Pupil struct {
	BaseModel
	Name      string `json:"name" csv:"姓名"`
	Class     Class  `json:"class" csv:"-"`
	ClassID   uint   `json:"-" csv:"班级ID"`
	CreatedBy string `json:"created_by" csv:"-"`
}
