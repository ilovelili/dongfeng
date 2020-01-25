package middleware

import (
	"fmt"
	"net/http"

	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Authenticator auth0 authenticator
type Authenticator struct {
	client *authingClient
}

// NewAuthenticator authenticator constructor
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		client: newAuthingClient(),
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
	return c.Path() == "/healthz" || c.Path() == "/*"
}

// TokenValidator jwt token validator
func (a *Authenticator) TokenValidator(accessToken string, c echo.Context) (bool, error) {
	userID, err := a.client.parseAccessToken(accessToken)
	if err != nil {
		return false, util.ResponseError(c, http.StatusUnauthorized, "401-100", "unauthorized", err)
	}

	userInfo, err := a.client.parseUserInfo(userID)
	if err != nil {
		return false, util.ResponseError(c, http.StatusUnauthorized, "401-101", "unauthorized", err)
	}

	// user name can't be empty
	if userInfo.Name == "" {
		return false, util.ResponseError(c, http.StatusUnauthorized, "401-102", "failed to parse user", err)
	}

	// add a email if email empty
	if userInfo.Email == "" {
		userInfo.Email = fmt.Sprintf("%s@dfyey.top", userID)
	}

	c.Set("userInfo", userInfo)
	return true, nil
}
