package model

// Class class entity
type Class struct {
	BaseModel
	Year      string `gorm:"unique_index:idx_year_name" json:"year" csv:"学年"`
	Name      string `gorm:"unique_index:idx_year_name" json:"name" csv:"班级"`
	CreatedBy string `json:"created_by" csv:"-"`
}
