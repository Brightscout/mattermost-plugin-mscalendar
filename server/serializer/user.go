package serializer

import (
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
	"golang.org/x/oauth2"
)

type User struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
	Mail              string `json:"mail,omitempty"`
}

type UserTokenHelpers struct {
	CheckUserStatus      func(logger bot.Logger, mattermostUserID string) bool
	ChangeUserStatus     func(err error, logger bot.Logger, mattermostUserID string, poster bot.Poster)
	RefreshAndStoreToken func(token *oauth2.Token, oconf *oauth2.Config, mattermostUserID string) (*oauth2.Token, error)
}
