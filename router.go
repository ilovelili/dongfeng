package main

import "github.com/ilovelili/dongfeng/handler"

// initRouter init router
func (a *App) initRouter() {
	s := a.Server
	s.Static("/", "")
	s.GET("/healthz", handler.Healthcheck)
	s.GET("/login", handler.Login)
	s.GET("/notifications", handler.Notifications)
	s.POST("/user/upload", handler.UploadAvatar)
	s.POST("/user/update", handler.UpdateUser)
}
