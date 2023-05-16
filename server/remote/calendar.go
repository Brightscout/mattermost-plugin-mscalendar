// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"time"

	"golang.org/x/oauth2"
)

type Calendar struct {
	Owner        *User   `json:"owner,omitempty"`
	ID           string  `json:"id"`
	Name         string  `json:"name,omitempty"`
	Events       []Event `json:"events,omitempty"`
	CalendarView []Event `json:"calendarView,omitempty"`
}

type ViewCalendarParams struct {
	StartTime    time.Time
	EndTime      time.Time
	RemoteUserID string
	AccessToken  *oauth2.Token
}

type ViewCalendarResponse struct {
	Error        *APIError
	RemoteUserID string
	Events       []*Event
}
