// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/kvstore"
	"github.com/pkg/errors"
)

const (
	ErrorUserInactive        = "You have been marked inactive because your refresh token is expired. Please disconnect and reconnect your account again."
	ErrorRefreshTokenExpired = "The refresh token has expired due to inactivity"
)

type UserStore interface {
	LoadUser(mattermostUserID string) (*User, error)
	LoadMattermostUserID(remoteUserID string) (string, error)
	LoadUserIndex() (UserIndex, error)
	StoreUser(user *User) error
	LoadUserFromIndex(mattermostUserID string) (*UserShort, error)
	DeleteUser(mattermostUserID string) error
	ModifyUserIndex(modify func(userIndex UserIndex) (UserIndex, error)) error
	StoreUserInIndex(user *User) error
	DeleteUserFromIndex(mattermostUserID string) error
	StoreUserActiveEvents(mattermostUserID string, events []string) error
}

type UserIndex []*UserShort

type UserShort struct {
	MattermostUserID string `json:"mm_id"`
	RemoteID         string `json:"remote_id"`
	Email            string `json:"email"`
}

type User struct {
	Settings          Settings `json:"mattermostSettings,omitempty"`
	Remote            *serializer.User
	OAuth2Token       *oauth2.Token
	PluginVersion     string
	MattermostUserID  string
	LastStatus        string
	WelcomeFlowStatus WelcomeFlowStatus `json:"mattermostFlags,omitempty"`
	ActiveEvents      []string          `json:"events"`
}

type Settings struct {
	DailySummary                      *DailySummaryUserSettings
	EventSubscriptionID               string
	UpdateStatus                      bool
	GetConfirmation                   bool
	ReceiveReminders                  bool
	ReceiveNotificationsDuringMeeting bool
}

type DailySummaryUserSettings struct {
	PostTime     string `json:"post_time"` // Kitchen format, i.e. 8:30AM
	Timezone     string `json:"tz"`        // Timezone in MSCal when PostTime is set/updated
	LastPostTime string `json:"last_post_time"`
	Enable       bool   `json:"enable"`
}

type WelcomeFlowStatus struct {
	PostIDs map[string]string
	Step    int
}

func (settings Settings) String() string {
	sub := "no subscription"
	if settings.EventSubscriptionID != "" {
		sub = "subscription ID: " + settings.EventSubscriptionID
	}
	return fmt.Sprintf(" - %s", sub)
}

func (user *User) Clone() *User {
	newUser := *user
	newRemoteUser := *user.Remote
	newUser.Remote = &newRemoteUser
	return &newUser
}

func (s *pluginStore) LoadUser(mattermostUserID string) (*User, error) {
	user := User{}
	err := kvstore.LoadJSON(s.userKV, mattermostUserID, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *pluginStore) LoadMattermostUserID(remoteUserID string) (string, error) {
	data, err := s.mattermostUserIDKV.Load(remoteUserID)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *pluginStore) LoadUserIndex() (UserIndex, error) {
	users := UserIndex{}
	err := kvstore.LoadJSON(s.userIndexKV, "", &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *pluginStore) LoadUserFromIndex(mattermostUserID string) (*UserShort, error) {
	users, err := s.LoadUserIndex()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.MattermostUserID == mattermostUserID {
			return u, nil
		}
	}

	return nil, ErrNotFound
}

func (s *pluginStore) StoreUser(user *User) error {
	err := kvstore.StoreJSON(s.userKV, user.MattermostUserID, user)
	if err != nil {
		return err
	}

	err = s.mattermostUserIDKV.Store(user.Remote.ID, []byte(user.MattermostUserID))
	if err != nil {
		_ = s.userKV.Delete(user.MattermostUserID)
		return err
	}
	return nil
}

func (s *pluginStore) DeleteUser(mattermostUserID string) error {
	u, err := s.LoadUser(mattermostUserID)
	if err != nil {
		return err
	}
	err = s.userKV.Delete(mattermostUserID)
	if err != nil {
		return err
	}
	err = s.mattermostUserIDKV.Delete(u.Remote.ID)
	if err != nil {
		return err
	}

	var userIndex []*UserShort
	err = kvstore.LoadJSON(s.userIndexKV, "", &userIndex)
	if err != nil {
		return err
	}
	filtered := []*UserShort{}
	for _, u := range userIndex {
		if u.MattermostUserID != mattermostUserID {
			filtered = append(filtered, u)
		}
	}
	err = kvstore.StoreJSON(s.userIndexKV, "", &filtered)
	if err != nil {
		return err
	}

	return nil
}

func (s *pluginStore) ModifyUserIndex(modify func(userIndex UserIndex) (UserIndex, error)) error {
	return kvstore.AtomicModify(s.userIndexKV, "", func(initial []byte, storeErr error) ([]byte, error) {
		if storeErr != nil && storeErr != ErrNotFound {
			return initial, storeErr
		}

		var storedIndex UserIndex
		if len(initial) > 0 {
			err := json.Unmarshal(initial, &storedIndex)
			if err != nil {
				return nil, err
			}
		}

		updated, err := modify(storedIndex)
		if err != nil {
			return nil, err
		}

		result, err := json.Marshal(updated)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
}

func (s *pluginStore) StoreUserInIndex(user *User) error {
	return s.ModifyUserIndex(func(userIndex UserIndex) (UserIndex, error) {
		newUser := &UserShort{
			MattermostUserID: user.MattermostUserID,
			RemoteID:         user.Remote.ID,
			Email:            user.Remote.Mail,
		}

		for i, u := range userIndex {
			if u.MattermostUserID == user.MattermostUserID && u.RemoteID == user.Remote.ID {
				var result UserIndex
				result = append(result, userIndex[:i]...)
				result = append(result, newUser)

				return append(result, userIndex[i+1:]...), nil
			}
		}

		return append(userIndex, newUser), nil
	})
}

func (s *pluginStore) DeleteUserFromIndex(mattermostUserID string) error {
	return s.ModifyUserIndex(func(userIndex UserIndex) (UserIndex, error) {
		for i, u := range userIndex {
			if u.MattermostUserID == mattermostUserID {
				return append(userIndex[:i], userIndex[i+1:]...), nil
			}
		}
		return userIndex, nil
	})
}

func (s *pluginStore) StoreUserActiveEvents(mattermostUserID string, events []string) error {
	u, err := s.LoadUser(mattermostUserID)
	if err != nil {
		return err
	}
	u.ActiveEvents = events
	return kvstore.StoreJSON(s.userKV, mattermostUserID, u)
}

func (index UserIndex) ByMattermostID() map[string]*UserShort {
	result := map[string]*UserShort{}

	for _, u := range index {
		result[u.MattermostUserID] = u
	}

	return result
}

func (index UserIndex) ByRemoteID() map[string]*UserShort {
	result := map[string]*UserShort{}

	for _, u := range index {
		result[u.RemoteID] = u
	}

	return result
}

func (index UserIndex) ByEmail() map[string]*UserShort {
	result := map[string]*UserShort{}

	for _, u := range index {
		result[u.Email] = u
	}

	return result
}

func (index UserIndex) GetMattermostUserIDs() []string {
	result := []string{}

	for _, u := range index {
		result = append(result, u.MattermostUserID)
	}

	return result
}

func GetCheckUserStatus(store Store, logger bot.Logger, mattermostUserID string) func() bool {
	return func() bool {
		user, err := store.LoadUser(mattermostUserID)
		if err != nil {
			logger.Errorf("Not able to load the user %s. %s", mattermostUserID, err.Error())
			return false
		}

		// Checking if the user is marked as inactive
		if user.OAuth2Token == nil {
			return false
		}

		return true
	}
}

func GetChangeUserStatus(store Store, logger bot.Logger, mattermostUserID string, poster bot.Poster) func(error) {
	return func(err error) {
		if err == nil {
			return
		}

		if !strings.Contains(err.Error(), ErrorRefreshTokenExpired) {
			return
		}

		user, err := store.LoadUser(mattermostUserID)
		if err != nil {
			return
		}

		// Marking the user as inactive
		user.OAuth2Token = nil
		if err = store.StoreUser(user); err != nil {
			return
		}

		if _, err := poster.DM(mattermostUserID, ErrorUserInactive); err != nil {
			logger.Errorf("Not able to DM the user %s. %s", mattermostUserID, err.Error())
		}
	}
}

// refreshAndStoreToken checks whether the current access token is expired or not. If it is,
// then it refreshes the token and stores the new pair of access and refresh tokens in kv store.
func RefreshAndStoreToken(store Store, token *oauth2.Token, oconf *oauth2.Config, mattermostUserID string) (*oauth2.Token, error) {
	// If there is only five minute left for the token to expire, we are refreshing the token.
	// We don't want the token to expire between the time when we decide that the old token is valid
	// and the time at which we create the request. We are handling that by not letting the token expire.
	if time.Until(token.Expiry) > 5*time.Minute {
		return token, nil
	}

	src := oconf.TokenSource(context.Background(), token)
	newToken, err := src.Token() // this actually goes and renews the tokens
	if err != nil {
		return nil, errors.Wrap(err, "unable to get the new refreshed token")
	}
	if newToken.AccessToken != token.AccessToken {
		user, err := store.LoadUser(mattermostUserID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get the new refreshed token")
		}
		user.OAuth2Token = newToken

		if err := store.StoreUser(user); err != nil {
			return nil, err
		}

		return newToken, nil
	}

	return token, nil
}
