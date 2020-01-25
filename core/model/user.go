package model

// User  user entity
type User struct {
	BaseModel
	Email   string `gorm:"unique_index" json:"email"`
	Name    string `json:"name"`
	Photo   string `gorm:"size:1024" json:"photo"`
	Enabled bool   `gorm:"default:true" json:"enabled"`
}
