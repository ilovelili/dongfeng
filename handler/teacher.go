package handler

import (
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetTeachers GET /teachers
func GetTeachers(c echo.Context) error {
	class, year := c.QueryParam("class"), c.QueryParam("year")
	teachers, err := teacherRepo.Find(class, year)
	if err != nil {
		return util.ResponseError(c, "500-110", "failed to get teachers", err)
	}

	return c.JSON(http.StatusOK, teachers)
}

// UpdateTeacher PUT /teacher
func UpdateTeacher(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	teacher := new(model.Teacher)
	if err := c.Bind(teacher); err != nil {
		return util.ResponseError(c, "400-107", "failed to bind teacher", err)
	}

	if teacher.Email != nil {
		user, err := userRepo.FindByEmail(*teacher.Email)
		// found
		if err == nil {
			teacher.User = user
			teacher.UserID = &user.ID
		} else {
			teacher.User = nil
			teacher.UserID = nil
		}
	}

	teacher.CreatedBy = userInfo.Email
	if err := teacherRepo.Save(teacher); err != nil {
		return util.ResponseError(c, "500-111", "failed to save teachers", err)
	}

	notify(model.TeacherListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// SaveTeachers POST /teachers
func SaveTeachers(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, "400-106", "failed to parse teachers", err)
	}
	defer file.Close()

	teachers := []*model.Teacher{}
	if err := gocsv.Unmarshal(file, &teachers); err != nil {
		return util.ResponseError(c, "400-106", "failed to parse teachers", err)
	}

	for _, teacher := range teachers {
		if teacher.Email != nil {
			user, err := userRepo.FindByEmail(*teacher.Email)
			// found
			if err == nil && user != nil {
				teacher.User = user
				teacher.UserID = &user.ID
			}
		}

		teacher.CreatedBy = userInfo.Email
	}

	if err := teacherRepo.DeleteInsert(teachers); err != nil {
		return util.ResponseError(c, "500-111", "failed to save teachers", err)
	}

	notify(model.TeacherListUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
