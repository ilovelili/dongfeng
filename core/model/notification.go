package model

import (
	"fmt"
	"strconv"
)

// Notification notification entity
type Notification struct {
	BaseModel
	User       string       `gorm:"index" json:"user"`
	CustomCode string       `json:"custom_code"`
	Details    string       `json:"details"`
	Link       *string      `gorm:"size:1024" json:"link"`
	Category   CategoryEnum `json:"category_id"`
	Read       bool         `gorm:"default:false" json:"read"`
}

// newNotification constructor
func newNotification(email, customcode, details string, link ...string) *Notification {
	notification := &Notification{
		User:       email,
		CustomCode: customcode,
		Details:    details,
	}

	if len(link) > 0 {
		notification.Link = &link[0]
	}

	notification.Category = resolveCategory(customcode)
	return notification
}

// System
var (
	// Broadcast system broadcast
	Broadcast = func(title, content string, link ...string) *Notification {
		return newNotification("AgentSmith", "N5001", fmt.Sprintf(`{"title":"%s","content":"%s"}`, title, content), link...)
	}
)

// GrowthProfile
var (
	// ProfileUpdated profile updated
	ProfileUpdated = func(email string) *Notification {
		return newNotification(email, "N1001", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// PupilListUpdated namelist updated
	PupilListUpdated = func(email string) *Notification {
		return newNotification(email, "N1002", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// ClassListUpdated classlist updated
	ClassListUpdated = func(email string) *Notification {
		return newNotification(email, "N1003", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// TeacherListUpdated teacher list updated
	TeacherListUpdated = func(email string) *Notification {
		return newNotification(email, "N1004", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// AttendanceUpdated attendance updated
	AttendanceUpdated = func(email string) *Notification {
		return newNotification(email, "N1005", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// EbookUpdated ebook updated
	EbookUpdated = func(email string) *Notification {
		return newNotification(email, "N1006", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// GrowthProfileUpdated growth profile template updated
	GrowthProfileUpdated = func(email string) *Notification {
		return newNotification(email, "N1007", fmt.Sprintf(`{"id":"%s"}`, email))
	}

	// GrowthProfileTemplateUpdated growth profile template updated
	GrowthProfileTemplateUpdated = func(email string) *Notification {
		return newNotification(email, "N1008", fmt.Sprintf(`{"id":"%s"}`, email))
	}
)

// Physique
var (
	PhysiqueUpdated = func(email string) *Notification {
		return newNotification(email, "N2001", fmt.Sprintf(`{"id":"%s"}`, email))
	}
)

// Nutrition
var (
	// RecipeUpdated recipe updated
	RecipeUpdated = func(email string) *Notification {
		return newNotification(email, "N3001", fmt.Sprintf(`{"id":"%s"}`, email))
	}
	// MenuUpdated menu updated by AgentSmith
	MenuUpdated = func() *Notification {
		return newNotification("AgentSmith", "N3002", "")
	}
	// IngredientUpdated ingredient updated
	IngredientUpdated = func(email string) *Notification {
		return newNotification(email, "N3003", fmt.Sprintf(`{"id":"%s"}`, email))
	}
)

// resolveCategory resolve category by custom code
func resolveCategory(customCode string) CategoryEnum {
	result, _ := strconv.Atoi(customCode[1:2])
	return CategoryEnum(result)
}
