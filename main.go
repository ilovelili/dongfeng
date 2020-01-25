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

	app.initConfig()
	app.initServer()
	app.initRouter()

	logger.Type("application").WithFields(logger.Fields{"host": host}).Infoln("starts dongfeng server")
	app.Server.Logger.Fatal(app.Server.Start(host))
}

func (a *App) initConfig() {
	a.Config = util.LoadConfig()
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

	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			req := c.Request()
			res := c.Response()
			if req.URL.Path != "/healthz" || res.Status != 200 {
				logger.Type("access").WithFields(logger.Fields{
					"remote_ip":       c.RealIP(),
					"host":            req.Host,
					"uri":             req.RequestURI,
					"method":          req.Method,
					"path":            req.URL.Path,
					"user_agent":      req.UserAgent(),
					"request-body":    string(reqBody),
					"response-status": res.Status,
					"response-body":   string(resBody),
				}).Infoln("")
			}
		},
	}))

	auth := mw.NewAuthenticator()
	e.Use(auth.Middleware())
	a.Server = e
}

// initRouter init router
func (a *App) initRouter() {
	s := a.Server
	s.Static("/", "")
	s.GET("/healthz", handler.Healthcheck)
	s.GET("/login", handler.Login)
}

func main() {
	flag.Parse()
	logger.SetAppName(appName)
	host := fmt.Sprintf(":%s", *port)
	Bootstrap(host)
}
