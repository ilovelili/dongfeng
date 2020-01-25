package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	jose "gopkg.in/square/go-jose.v2"
)


// Authenticator auth0 authenticator
type Authenticator struct {
	client *AuthingClient
}

// NewAuthenticator authenticator constructor
func NewAuthenticator() *Authenticator {		
	return &Authenticator{
		client: NewAuthingClient(),
	}
}

// Middleware authenticator middleware func
func (a *Authenticator) Middleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Skipper:    a.Skipper,
		Validator:  a.TokenValidator,
	})
}

// Skipper auth skipper
func (a *Authenticator) Skipper(c echo.Context) bool {
	return c.Path() == "/healthz"
}

// TokenValidator jwt token validator
func (a *Authenticator) TokenValidator(idtoken string, c echo.Context) (bool, error) {
	status, err := a.client.verifyLogin(idtoken)
	if !status || err != nil {
		return false, util.ResponseError(c, http.StatusUnauthorized, "401-100", "unauthorized", err)
	}
	
	userinfo, err := a.client.parseUserInfo(pid)
	if err != nil {
		return false, util.ResponseError(c, http.StatusUnauthorized, "401-100", "unauthorized", err)
	}
	

	c.Set(util.ContextUser, claims)
	return true, nil
}
