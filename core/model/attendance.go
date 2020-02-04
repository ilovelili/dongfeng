package model

import "github.com/ilovelili/dongfeng/util"

// Absence absence entity
type Absence struct {
	BaseModel
	CSVYear   string `gorm:"-" csv:"学年"`
	Date      string `csv:"日期"`
	CSVClass  string `gorm:"-" csv:"班级"`
	CSVName   string `gorm:"-" csv:"姓名"`
	Pupil     *Pupil `csv:"-"`
	PupilID   uint   `csv:"-"`
	CreatedBy string `csv:"-"`
}

// Attendance attendence entity
type Attendance struct {
	Date        string          `json:"date"`
	HolidayType HolidayTypeEnum `json:"holiday"`
	Pupil       *Pupil          `json:"pupil"`
	Absent      bool            `json:"absent"` // attended or absent
}

// Holiday holiday entity
type Holiday struct {
	ID          uint `gorm:"primary_key" json:"id"`
	From        string
	To          string
	Description string
	CreatedBy   string
}

// IsHoliday date is holiday
func (h *Holiday) IsHoliday(d string) bool {
	if d == h.From || d == h.To {
		return true
	}
	fromDate, _ := util.StringToTime(h.From)
	toDate, _ := util.StringToTime(h.To)
	date, _ := util.StringToTime(d)
	return date.After(fromDate) && date.Before(toDate)
}

// HolidayType holiday type
type HolidayType struct {
	Date string          `json:"date"`
	Type HolidayTypeEnum `json:"type"`
}

// HolidayTypeEnum 0: working day | 1: weekend | 2: holiday
type HolidayTypeEnum uint

const (
	// WorkingDays working days
	WorkingDays HolidayTypeEnum = 0
	// Weekends weekends
	Weekends HolidayTypeEnum = 1
	// Holidays holidays
	Holidays HolidayTypeEnum = 2
)
