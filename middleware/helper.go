package middleware

import (
	"encoding/json"
	"sync"

	"fmt"

	authing "github.com/Authing/authing-go-sdk"
	"github.com/go-resty/resty/v2"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/ilovelili/logger"
	"github.com/kelvinji2009/graphql"
)

var (
	once     sync.Once
	instance *authingClient
)

// authingClient authing client wrapper
type authingClient struct {
	client   *authing.Client
	clientID string
}

// ParseTokenResponse parse token response
type ParseTokenResponse struct {
	Status bool  `json:"status"`
	Code   int   `json:"code"`
	Token  Token `json:"token,omitempty"`
}

// Token token data
type Token struct {
	Data `json:"data"`
}

// Data token data fields
type Data struct {
	UnionID  string `json:"unionid"`
	ID       string `json:"id"`
	ClientID string `json:"clientId"`
}

// newAuthingClient new authing client
func newAuthingClient() *authingClient {
	config := util.LoadConfig()
	clientID, appSecret := config.Auth.ClientID, config.Auth.ClientSecret
	once.Do(func() {
		client := authing.NewClient(clientID, appSecret, false)
		client.Client.Log = func(s string) {
			logger.Type("authing").Infoln(s)
		}
		instance = &authingClient{
			client:   client,
			clientID: clientID,
		}
	})

	return instance
}

// parseAccessToken parse access token
func (c *authingClient) parseAccessToken(accessToken string) (userID string, err error) {
	client := resty.New()
	rawResp, err := client.R().
		SetQueryParams(map[string]string{"access_token": accessToken}).
		SetHeader("Accept", "application/json").
		Get("https://users.authing.cn/authing/token")
	if err != nil {
		return
	}

	var resp ParseTokenResponse
	err = json.Unmarshal(rawResp.Body(), &resp)
	if err == nil {
		userID = resp.Token.ID
	}
	return
}

// parseUserInfo parse user info based on user id
func (c *authingClient) parseUserInfo(userID string) (model.User, error) {
	user := model.User{}
	p := authing.UserQueryParameter{
		ID:               graphql.String(userID),
		RegisterInClient: graphql.String(c.clientID),
	}

	q, err := c.client.User(&p)
	if err != nil {
		return user, err
	}

	authingUser := q.User
	if authingUser.Blocked {
		return user, fmt.Errorf("This user is blocked")
	}
	if authingUser.IsDeleted {
		return user, fmt.Errorf("This user has been deleted")
	}

	user.Email = string(authingUser.Email)
	user.Name = string(authingUser.Username)
	user.Photo = string(authingUser.Photo)
	return user, nil
}
