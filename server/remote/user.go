// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"golang.org/x/oauth2"
)

type WorkingHours struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TimeZone  struct {
		Name string `json:"name"`
	}
	DaysOfWeek []string `json:"daysOfWeek"`
}

type MailboxSettings struct {
	TimeZone     string       `json:"timeZone"`
	WorkingHours WorkingHours `json:"workingHours"`
}

type UserTokenHelpers struct {
	CheckUserStatus      func(mattermostUserID string) bool
	ChangeUserStatus     func(err error, mattermostUserID string)
	RefreshAndStoreToken func(token *oauth2.Token, oconf *oauth2.Config, mattermostUserID string) (*oauth2.Token, error)
}
