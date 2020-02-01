package main

import "github.com/ilovelili/dongfeng/handler"

// initRouter init router
func (a *App) initRouter() {
	s := a.Server
	s.Static("/", "")
	s.GET("/healthz", handler.Healthcheck)
	s.GET("/login", handler.Login)

	s.GET("/notifications", handler.GetNotifications)
	s.POST("/notifications", handler.SetNotificationsRead)

	s.POST("/user/upload", handler.UploadAvatar)
	s.PUT("/user/update", handler.UpdateUser)

	s.GET("/classes", handler.GetClasses)
	s.POST("/classes", handler.SaveClasses)

	s.GET("/pupils", handler.GetPupils)
	s.POST("/pupils", handler.SavePupils)

	s.GET("/teachers", handler.GetTeachers)
	s.PUT("/teacher", handler.UpdateTeacher)
	s.POST("/teachers", handler.SaveTeachers)

	s.GET("/attendances", handler.GetAttendances)
	s.PUT("/attendance", handler.UpdateAttendance)
	s.POST("/attendances", handler.SaveAbsences)

	s.GET("/physiques", handler.GetPhysiques)
	s.PUT("/physique", handler.UpdatePhysique)
	s.POST("/physiques", handler.SavePhysiques)

	s.GET("/ingredients", handler.GetIngredients)
	s.POST("/ingredients", handler.SaveIngredients)

	s.GET("/masters", handler.GetMasters)

}
