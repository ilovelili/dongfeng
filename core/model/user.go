package model

// User  user entity
type User struct {
	BaseModel
	Email string   `gorm:"unique_index" json:"email"`
	Name  string   `json:"name"`
	Photo string   `gorm:"size:1024" json:"photo"`
	Role  RoleEnum `gorm:"default:5" json:"role"`
}

// RoleEnum roles enum
type RoleEnum uint

const (
	// RoleUndefined undefined
	RoleUndefined RoleEnum = 0
	// RoleAgentSmith super user
	RoleAgentSmith RoleEnum = 1
	// RoleAdmin admin role
	RoleAdmin RoleEnum = 2
	// RoleTeacher teacher role
	RoleTeacher RoleEnum = 3
	// RoleHealth role health
	RoleHealth RoleEnum = 4
	// RoleNormal normal role
	RoleNormal RoleEnum = 5
	// RoleOthers others role
	RoleOthers RoleEnum = 6
)
