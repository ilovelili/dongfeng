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
	s.POST("/user/update", handler.UpdateUser)
	s.GET("/classes", handler.GetClasses)
	s.POST("/classes", handler.UpdateClasses)
	s.GET("/pupils", handler.GetPupils)
	s.POST("/pupils", handler.UpdatePupils)
	s.GET("/teachers", handler.GetTeachers)
	s.POST("/teacher", handler.UpdateTeacher)
	s.POST("/teachers", handler.UpdateTeachers)
}
