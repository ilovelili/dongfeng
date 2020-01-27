package model

// CategoryEnum category enum
type CategoryEnum string

const (
	// CategoryReserved reserved
	CategoryReserved CategoryEnum = ""
	// CategoryGrowthProfile growth profile
	CategoryGrowthProfile CategoryEnum = "成长档案"
	// CategoryPhysique physique
	CategoryPhysique CategoryEnum = "体格信息"
	// CategoryNutrition nutrition
	CategoryNutrition CategoryEnum = "营养膳食"
	// CategoryAttendance attendance
	CategoryAttendance CategoryEnum = "出勤信息"
	// CategoryAgentSmith agent smith
	CategoryAgentSmith CategoryEnum = "系统信息"
)
