package serializer

import "golang.org/x/oauth2"

type User struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
	Mail              string `json:"mail,omitempty"`
}

type UserTokenHelpers struct {
	CheckUserStatus      func() bool
	ChangeUserStatus     func(error)
	RefreshAndStoreToken func(token *oauth2.Token, oconf *oauth2.Config, mattermostUserID string) (*oauth2.Token, error)
}
