// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"time"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
)

type Calendar struct {
	Owner        *serializer.User   `json:"owner,omitempty"`
	ID           string             `json:"id"`
	Name         string             `json:"name,omitempty"`
	Events       []serializer.Event `json:"events,omitempty"`
	CalendarView []serializer.Event `json:"calendarView,omitempty"`
}

type ViewCalendarParams struct {
	StartTime    time.Time
	EndTime      time.Time
	RemoteUserID string
}

type ViewCalendarResponse struct {
	Error        *APIError
	RemoteUserID string
	Events       []*serializer.Event
}
