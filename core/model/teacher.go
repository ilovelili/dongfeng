package model

// Teacher teacher entity
type Teacher struct {
	BaseModel
	Name      string  `json:"name" csv:"姓名"`
	Email     *string `json:"email" csv:"邮箱,omitempty"`
	User      *User   `json:"user" csv:"-"`
	UserID    *uint   `json:"-" csv:"-"`
	Class     *Class  `json:"class" csv:"-"`
	ClassID   *uint   `json:"-" csv:"班级ID,omitempty"`
	CreatedBy string  `json:"created_by" csv:"-"`
}
