package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ilovelili/dongfeng/handler"
	mw "github.com/ilovelili/dongfeng/middleware"
	"github.com/ilovelili/dongfeng/util"
	"github.com/ilovelili/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	port = flag.String("port", "8080", "port to listen")
)

const (
	appName = "dongfeng"
)

// App app
type App struct {
	Server *echo.Echo
	Config *util.Config
}

// Bootstrap bootstrap server
func Bootstrap(host string) {
	app := new(App)

	app.initServer()
	app.initRouter()

	logger.Type("application").WithFields(logger.Fields{"host": host}).Infoln("starts dongfeng server")
	app.Server.Logger.Fatal(app.Server.Start(host))
}

func (a *App) initServer() {
	e := echo.New()
	e.HideBanner = true

	aos := []string{"*"}
	if os.Getenv("DF_ALLOW_ORIGINS") != "" {
		aos = strings.Split(os.Getenv("DF_ALLOW_ORIGINS"), ",")
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: aos,
	}))

	// e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
	// 	Handler: func(c echo.Context, reqBody, resBody []byte) {
	// 		req := c.Request()
	// 		res := c.Response()
	// 		if (req.URL.Path != "/api/healthz" &&
	// 			req.URL.Path != "/api/profileContent" &&
	// 			req.URL.Path != "/api/profileTemplate") ||
	// 			res.Status != 200 {
	// 			logger.Type("access").WithFields(logger.Fields{
	// 				"remote_ip":       c.RealIP(),
	// 				"host":            req.Host,
	// 				"uri":             req.RequestURI,
	// 				"method":          req.Method,
	// 				"path":            req.URL.Path,
	// 				"user_agent":      req.UserAgent(),
	// 				"request-body":    string(reqBody),
	// 				"response-status": res.Status,
	// 				"response-body":   string(resBody),
	// 			}).Infoln("")
	// 		}
	// 	},
	// }))

	a.Server = e
}

// initRouter init router
func (a *App) initRouter() {
	s := a.Server

	// add "api" prefix for easier load balancer settings
	auth := mw.NewAuthenticator()
	api := s.Group("/api", auth.Middleware())

	api.File("/const.json", "const.json")

	api.GET("/healthz", handler.Healthcheck)
	api.GET("/login", handler.Login)

	api.GET("/notifications", handler.GetNotifications)
	api.POST("/notifications", handler.SetNotificationsRead)

	api.GET("/users", handler.GetUsers)
	api.POST("/user/upload", handler.UploadAvatar)
	api.PUT("/user/update", handler.UpdateUser)

	api.GET("/classes", handler.GetClasses)
	api.POST("/classes", handler.SaveClasses)

	api.GET("/pupils", handler.GetPupils)
	api.POST("/pupils", handler.SavePupils)

	api.GET("/teachers", handler.GetTeachers)
	api.PUT("/teacher", handler.UpdateTeacher)
	api.POST("/teachers", handler.SaveTeachers)

	api.GET("/attendances", handler.GetAttendances)
	api.PUT("/attendance", handler.UpdateAttendance)
	api.POST("/attendances", handler.SaveAbsences)

	api.GET("/physiques", handler.GetPhysiques)
	api.PUT("/physique", handler.UpdatePhysique)
	api.POST("/physiques", handler.SavePhysiques)

	api.GET("/menus", handler.GetMenus)
	api.GET("/recipes", handler.GetRecipes)
	api.GET("/ingredients", handler.GetIngredients)
	api.POST("/ingredients", handler.SaveIngredients)

	api.GET("/profileTemplateContent", handler.GetProfileTemplateContent)
	api.POST("/profileTemplateTags", handler.SaveProfileTemplateTags)
	api.POST("/profileTemplate", handler.SaveProfileTemplate)
	api.DELETE("/profileTemplate", handler.DeleteProfileTemplate)
	api.GET("/profileTemplates", handler.GetProfileTemplates)

	api.GET("/profiles", handler.GetProfiles)
	api.GET("/profile/prev", handler.GetPreviousProfile)
	api.GET("/profile/next", handler.GetNextProfile)
	api.POST("/profile", handler.SaveProfile)
	api.DELETE("/profile", handler.DeleteProfile)
	api.GET("/profileContent", handler.GetProfileContent)
	api.POST("/profileContent", handler.SaveProfileContent)

	api.GET("/ebooks", handler.GetEbooks)
	api.POST("/ebook", handler.UpdateEbook)

	api.GET("/masters", handler.GetMasters)
}

func main() {
	flag.Parse()
	logger.SetAppName(appName)
	host := fmt.Sprintf(":%s", *port)
	Bootstrap(host)
}
