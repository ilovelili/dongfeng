package model

// Notification notification entity
type Notification struct {
	BaseModel
	User       string       `gorm:"index" json:"user"`
	CustomCode string       `json:"custom_code"`
	Details    string       `json:"details"`
	Link       string       `gorm:"size:1024" json:"link"`
	Category   CategoryEnum `json:"category_id"`
	Read       bool         `json:"read"`
}
