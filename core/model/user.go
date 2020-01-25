package model

// User  user entity
type User struct {
	BaseModel
	Email   string   `gorm:"unique_index" json:"email"`
	Name    string   `json:"name"`
	Photo   string   `gorm:"size:1024" json:"photo"`
	Enabled bool     `gorm:"default:true" json:"enabled"`
	Role    RoleEnum `gorm:"default:1" json:"role"`
}

// RoleEnum roles enum
type RoleEnum uint

const (
	// RoleAgentSmith super user
	RoleAgentSmith RoleEnum = 0
	// RoleAdmin admin role
	RoleAdmin RoleEnum = 1
	// RoleNormal normal role
	RoleNormal RoleEnum = 2
	// RoleHealthcare role healthcare
	RoleHealthcare RoleEnum = 3
)
