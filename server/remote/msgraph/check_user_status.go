package msgraph

import "strings"

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
	if user.OAuth2Token == nil {
		return false
	}

	return true
}

func (c *client) ChangeUserStatus(err error) {
	if err == nil {
		return
	}

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
	user.OAuth2Token = nil
	if err = c.store.StoreUser(user); err != nil {
		c.Logger.Errorf("Not able to store the user %s. %s", c.mattermostUserID, err.Error())
		return
	}

	if _, err := c.poster.DM(c.mattermostUserID, ErrorUserInactive); err != nil {
		c.Logger.Errorf("Not able to DM the user %s. %s", c.mattermostUserID, err.Error())
	}
}
