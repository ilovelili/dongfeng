package middleware

import (
	"encoding/json"
	"sync"

	"fmt"

	authing "github.com/Authing/authing-go-sdk"
	"github.com/ilovelili/dongfeng/util"
	"github.com/kelvinji2009/graphql"
)

var (
	once     sync.Once
	instance *AuthingClient
)

// authingClient authing client wrapper
type authingClient struct {
	client   *authing.Client
	clientID string
}

// NewAuthClient  new authing client
func newAuthClient() *authingClient {
	config := util.LoadConfig()
	clientID, appSecret := config.Auth.ClientID, config.Auth.ClientSecret
	once.Do(func() {
		client := authing.NewClient(clientID, appSecret, false)
		instance = &authingClient{
			client:   client,
			clientID: clientID,
		}
	})

	return instance
}

// parseUserInfo parse user info based on user id
func (c *authingClient) parseUserInfo(userID string) (*model.User, error) {
	user := new(model.User)
	p := authing.UserQueryParameter{
		ID:               graphql.String(userID),
		RegisterInClient: graphql.String(c.ClientID),
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

	user.Email = authingUser.Email
	user.Name = authingUser.Username
	user.Picture = authingUser.Photo
	return user, nil
}