package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetAttendances GET /attendances
func GetAttendances(c echo.Context) error {
	var (
		err         error
		pupilID     uint = 0
		classID     uint = 0
		absences         = []*model.Absence{}
		pupils           = []*model.Pupil{}
		attendances      = []*model.Attendance{}
		year             = c.QueryParam("year")
		class            = c.QueryParam("class")
		pupil            = c.QueryParam("name")
		from             = c.QueryParam("from")
		to               = c.QueryParam("to")
	)

	if from == "" {
		from = util.TimeToString(time.Now().AddDate(0, -1, 0))
	}
	if to == "" {
		to = util.TimeToString(time.Now())
	}
	if class != "" {
		_classID, err := strconv.ParseUint(class, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-108", "invalid class", err)
		}
		classID = uint(_classID)
	}
	if pupil != "" {
		_pupilID, err := strconv.ParseUint(pupil, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-109", "invalid pupil", err)
		}
		pupilID = uint(_pupilID)
	}

	// retrieve pupils
	if pupilID != 0 {
		pupils, err = pupilRepo.FindByPupilID(pupilID)
	} else if classID != 0 {
		pupils, err = pupilRepo.FindByClassID(classID)
	} else {
		pupils, err = pupilRepo.Find("", year)
	}
	if err != nil {
		return util.ResponseError(c, "500-107", "failed to get pupils", err)
	}

	// retrieve absences
	if pupilID != 0 {
		absences, err = attendanceRepo.FindAbsencesByPupil(pupilID, from, to)
	} else if classID != 0 {
		absences, err = attendanceRepo.FindAbsencesByClass(classID, from, to)
	} else {
		absences, err = attendanceRepo.FindAbsences(from, to)
	}
	if err != nil {
		return util.ResponseError(c, "500-112", "failed to get absences", err)
	}

	// retrieve holidays
	holidays, err := attendanceRepo.FindHolidays()
	if err != nil {
		return util.ResponseError(c, "500-114", "failed to get holidays", err)
	}

	fromTime, _ := util.StringToTime(from)
	toTime, _ := util.StringToTime(to)

DAILYLOOP:
	for d := fromTime; d.Before(toTime.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		currentDate := util.TimeToString(d)

		// skip if it's weekend or holiday
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			attendances = append(attendances, &model.Attendance{
				Date:        currentDate,
				HolidayType: model.Weekends,
			})
			continue DAILYLOOP
		}

		for _, holiday := range holidays {
			if holiday.IsHoliday(currentDate) {
				attendances = append(attendances, &model.Attendance{
					Date:        currentDate,
					HolidayType: model.Holidays,
				})
				continue DAILYLOOP
			}
		}

		absentPupilList := []*model.Pupil{}
		for _, absence := range absences {
			if absence.Date == currentDate {
				absentPupilList = append(absentPupilList, absence.Pupil)
			}
		}

		for _, pupil := range pupils {
			attendance := &model.Attendance{
				Date:        currentDate,
				HolidayType: model.WorkingDays,
				Pupil:       pupil,
				Absent:      false,
			}

			// if pupil in absent list, set absent to true
			for _, absentPupil := range absentPupilList {
				if absentPupil.ID == pupil.ID {
					attendance.Absent = true
				}
			}

			attendances = append(attendances, attendance)
		}
	}

	return c.JSON(http.StatusOK, attendances)
}

// UpdateAttendance PUT /attendances
func UpdateAttendance(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)

	attendanceReq := new(struct {
		Pupil  uint   `json:"pupil"`
		Date   string `json:"date"`
		Absent bool   `json:"absent"`
	})

	if err := c.Bind(attendanceReq); err != nil {
		return util.ResponseError(c, "400-110", "failed to bind attendance", err)
	}

	absence := &model.Absence{
		PupilID:   attendanceReq.Pupil,
		Date:      attendanceReq.Date,
		CreatedBy: userInfo.Email,
	}

	if attendanceReq.Absent {
		if err := attendanceRepo.Save(absence); err != nil {
			return util.ResponseError(c, "500-113", "failed to save absence", err)
		}
	} else {
		if err := attendanceRepo.Delete(absence); err != nil {
			return util.ResponseError(c, "500-113", "failed to save absence", err)
		}
	}

	notify(model.AttendanceUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// SaveAbsences POST /attendances
func SaveAbsences(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, "400-111", "failed to parse absences", err)
	}
	defer file.Close()

	absences := []*model.Absence{}
	if err := gocsv.Unmarshal(file, &absences); err != nil {
		return util.ResponseError(c, "400-111", "failed to parse absences", err)
	}

	// year | classname | pupilname | pupil id
	classPupilMap := make(map[string]map[string]map[string]uint)
	for _, absence := range absences {
		year := absence.CSVYear
		if _, ok := classPupilMap[year]; !ok {
			classes, err := classRepo.Find(year)
			if err != nil {
				return util.ResponseError(c, "500-105", "failed to get classes", err)
			}

			classPupilMap[year] = make(map[string]map[string]uint)
			for _, class := range classes {
				pupils, err := pupilRepo.FindByClassID(class.ID)
				if err != nil {
					return util.ResponseError(c, "500-107", "failed to get pupils", err)
				}

				classPupilMap[year][class.Name] = make(map[string]uint)
				for _, pupil := range pupils {
					classPupilMap[year][class.Name][pupil.Name] = pupil.ID
				}
			}
		}

		found := false
		if map1, ok := classPupilMap[absence.CSVYear]; ok {
			if map2, ok := map1[absence.CSVClass]; ok {
				if id, ok := map2[absence.CSVName]; ok {
					absence.PupilID = id
					found = true
				}
			}
		}

		if !found {
			return util.ResponseError(c, "400-109", "invalid pupil", err)
		}

		absence.CreatedBy = userInfo.Email
	}

	if err := attendanceRepo.SaveAll(absences); err != nil {
		return util.ResponseError(c, "500-113", "failed to save absence", err)
	}

	notify(model.AttendanceUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
