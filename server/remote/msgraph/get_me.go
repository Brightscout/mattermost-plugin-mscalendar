// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
	"github.com/pkg/errors"
)

const (
	ErrorUserInActive        = "You have been marked inactive. Please disconnect and reconnect your account again."
	LogUserInActive          = "User %s is inactive. Please disconnect and reconnect your account again."
	ErrorRefreshTokenExpired = "The refresh token has expired due to inactivity"
)

func (c *client) GetMe() (*serializer.User, error) {
	graphUser, err := c.rbuilder.Me().Request().Get(c.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "msgraph GetMe")
	}

	if graphUser.ID == nil {
		return nil, errors.New("user has no ID")
	}
	if graphUser.DisplayName == nil {
		return nil, errors.New("user has no Display Name")
	}
	if graphUser.UserPrincipalName == nil {
		return nil, errors.New("user has no Principal Name")
	}
	if graphUser.Mail == nil {
		return nil, errors.New("user has no email address. Make sure the Microsoft account is associated to an Outlook product")
	}

	user := &serializer.User{
		ID:                *graphUser.ID,
		DisplayName:       *graphUser.DisplayName,
		UserPrincipalName: *graphUser.UserPrincipalName,
		Mail:              *graphUser.Mail,
	}

	return user, nil
}

func (c *client) CheckUserStatus() bool {
	// Not checking for API calls made using super client
	if c.store == nil {
		return true
	}

	user, err := c.store.LoadUser(c.mattermostUserID)
	if err != nil {
		c.Logger.Errorf("Not able to load the user %s. %s", c.mattermostUserID, err.Error())
		return false
	}

	// Checking if the user is marked as inactive
	if user.OAuth2Token.AccessToken == "" {
		if _, err := c.poster.DM(c.mattermostUserID, ErrorUserInActive); err != nil {
			c.Logger.Errorf("Not able to DM the user %s. %s", c.mattermostUserID, err.Error())
		}
		return false
	}

	return true
}

func (c *client) ChangeUserStatus(err error) {
	if !strings.Contains(err.Error(), ErrorRefreshTokenExpired) {
		return
	}

	// Not checking for API calls made using super client
	if c.store == nil {
		return
	}

	user, err := c.store.LoadUser(c.mattermostUserID)
	if err != nil {
		c.Logger.Errorf("Not able to load the user %s. %s", c.mattermostUserID, err.Error())
		return
	}

	// Marking the user as inactive
	user.OAuth2Token.AccessToken = ""
	if c.store.StoreUser(user); err != nil {
		c.Logger.Errorf("Not able to store the user %s. %s", c.mattermostUserID, err.Error())
		return
	}
}
