// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

import (
	"net/url"
	"time"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
)

type Client interface {
	AcceptEvent(remoteUserID, eventID string) error
	CallFormPost(method, path string, in url.Values, out interface{}) (responseData []byte, err error)
	CallJSON(method, path string, in, out interface{}) (responseData []byte, err error)
	CreateCalendar(remoteUserID string, calendar *Calendar) (*Calendar, error)
	CreateEvent(remoteUserID string, calendarEvent *serializer.Event) (*serializer.Event, error)
	CreateMySubscription(notificationURL string) (*serializer.Subscription, error)
	DeclineEvent(remoteUserID, eventID string) error
	DeleteCalendar(remoteUserID, calendarID string) error
	DeleteSubscription(subscriptionID string) error
	FindMeetingTimes(remoteUserID string, meetingParams *FindMeetingTimesParameters) (*MeetingTimeSuggestionResults, error)
	GetCalendars(remoteUserID string) ([]*Calendar, error)
	GetDefaultCalendarView(remoteUserID string, startTime, endTime time.Time) ([]*serializer.Event, error)
	DoBatchViewCalendarRequests([]*ViewCalendarParams) ([]*ViewCalendarResponse, error)
	GetEvent(remoteUserID, eventID string) (*serializer.Event, error)
	GetMailboxSettings(remoteUserID string) (*MailboxSettings, error)
	GetMe() (*serializer.User, error)
	GetNotificationData(*Notification) (*Notification, error)
	GetSchedule(requests []*ScheduleUserInfo, startTime, endTime *serializer.DateTime, availabilityViewInterval int) ([]*ScheduleInformation, error)
	ListSubscriptions() ([]*serializer.Subscription, error)
	RenewSubscription(subscriptionID string) (*serializer.Subscription, error)
	TentativelyAcceptEvent(remoteUserID, eventID string) error
	GetSuperuserToken() (string, error)
}
