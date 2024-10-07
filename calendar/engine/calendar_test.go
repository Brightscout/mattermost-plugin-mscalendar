package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/config"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/engine/mock_plugin_api"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote/mock_remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store/mock_store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/utils/bot/mock_bot"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/require"
)

func TestViewCalendar(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	now := time.Now()
	from := now.Add(-time.Hour)
	to := now.Add(time.Hour)

	tests := []struct {
		name       string
		user       *User
		setupMock  func()
		assertions func(t *testing.T, events []*remote.Event, err error)
	}{
		{
			name: "error filtering with client",
			user: &User{
				User:             nil,
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockStore.EXPECT().LoadUser(gomock.Any()).Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
			},
		},
		{
			name: "error getting calendar view",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockClient.EXPECT().GetDefaultCalendarView("testRemoteUserID", from, to).Return(nil, fmt.Errorf("error getting calendar view")).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error getting calendar view")
			},
		},
		{
			name: "successful calendar view",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockClient.EXPECT().GetDefaultCalendarView("testRemoteUserID", from, to).Return([]*remote.Event{{Subject: "Test Event"}}, nil).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.NoError(t, err)
				require.NotNil(t, events)
				require.Len(t, events, 1)
				require.Equal(t, "Test Event", events[0].Subject, "Expected first event's subject to be %s, but got %s", "Test Event", events[0].Subject)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			events, err := mscalendar.ViewCalendar(tt.user, from, to)

			tt.assertions(t, events, err)
		})
	}
}

func TestGetTodayCalendarEvents(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	now := time.Now()
	timezone := "America/Los_Angeles"
	from, to := getTodayHoursForTimezone(now, timezone)

	tests := []struct {
		name       string
		user       *User
		setupMock  func()
		assertions func(t *testing.T, events []*remote.Event, err error)
	}{
		{
			name: "error expanding remote user",
			user: &User{
				User:             nil,
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
			},
		},
		{
			name: "error getting calendar view",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockClient.EXPECT().GetDefaultCalendarView("testRemoteUserID", from, to).Return(nil, fmt.Errorf("error getting calendar view")).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error getting calendar view")
			},
		},
		{
			name: "successful calendar view",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMID"},
				MattermostUserID: "testMMUserID",
			},
			setupMock: func() {
				mockClient.EXPECT().GetDefaultCalendarView("testRemoteUserID", from, to).Return([]*remote.Event{{Subject: "Today's Test Event"}}, nil).Times(1)
			},
			assertions: func(t *testing.T, events []*remote.Event, err error) {
				require.NoError(t, err)
				require.NotNil(t, events)
				require.Len(t, events, 1)
				require.Equal(t, "Today's Test Event", events[0].Subject, "Expected first event's subject to be %s, but got %s", "Today's Test Event", events[0].Subject)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			events, err := mscalendar.getTodayCalendarEvents(tt.user, now, timezone)

			tt.assertions(t, events, err)
		})
	}
}

func TestCreateCalendar(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	tests := []struct {
		name       string
		user       *User
		calendar   *remote.Calendar
		setupMock  func()
		assertions func(t *testing.T, createdCalendar *remote.Calendar, err error)
	}{
		{
			name: "error expanding user",
			user: &User{
				MattermostUserID: "testMMUserID",
			},
			calendar: &remote.Calendar{
				Name: "Test Calendar",
			},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, createdCalendar *remote.Calendar, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
			},
		},
		{
			name: "error creating calendar",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMUserID"},
				MattermostUserID: "testMMUserID",
			},
			calendar: &remote.Calendar{
				Name: "Test Calendar",
			},
			setupMock: func() {
				mockClient.EXPECT().CreateCalendar("testRemoteUserID", &remote.Calendar{Name: "Test Calendar"}).Return(nil, fmt.Errorf("error creating calendar")).Times(1)
			},
			assertions: func(t *testing.T, createdCalendar *remote.Calendar, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error creating calendar")
			},
		},
		{
			name: "successful calendar creation",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUser:   &model.User{Id: "testMMUserID"},
				MattermostUserID: "testMMUserID",
			},
			calendar: &remote.Calendar{
				Name: "Test Calendar",
			},
			setupMock: func() {
				mockClient.EXPECT().CreateCalendar("testRemoteUserID", &remote.Calendar{Name: "Test Calendar"}).Return(&remote.Calendar{Name: "Created Test Calendar"}, nil).Times(1)
			},
			assertions: func(t *testing.T, createdCalendar *remote.Calendar, err error) {
				require.NoError(t, err)
				require.NotNil(t, createdCalendar)
				require.Equal(t, "Created Test Calendar", createdCalendar.Name, "Expected calendar name to be %s, but got %s", "Created Test Calendar", createdCalendar.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			createdCalendar, err := mscalendar.CreateCalendar(tt.user, tt.calendar)
			tt.assertions(t, createdCalendar, err)
		})
	}
}

func TestCreateEvent(t *testing.T) {
	mscalendar, mockStore, mockPoster, _, mockPluginAPI, mockClient, mockLogger := MockSetup(t)

	tests := []struct {
		name              string
		user              *User
		event             *remote.Event
		mattermostUserIDs []string
		setupMock         func()
		assertions        func(t *testing.T, createdEvent *remote.Event, err error)
		expectedEvent     *remote.Event
	}{
		{
			name: "error expanding user",
			user: &User{
				MattermostUserID: "testMMUserID",
			},
			event:             &remote.Event{Subject: "Test Event"},
			mattermostUserIDs: []string{"testMMUserID1"},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, createdEvent *remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
			},
		},
		{
			name: "error creating direct message",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUserID: "testMMUserID",
			},
			event:             &remote.Event{Subject: "Test Event"},
			mattermostUserIDs: []string{"testMMUserID1"},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID1").Return(nil, errors.New("not found")).Times(1)
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockPoster.EXPECT().DM("testMMUserID1", gomock.Any(), "testDisplayName", "testDisplayName", "testCommandTrigger").Return("", fmt.Errorf("error creating DM")).Times(1)
				mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())
				mockClient.EXPECT().CreateEvent("testRemoteUserID", gomock.Any()).Return(&remote.Event{}, nil).Times(1)
			},
			assertions: func(t *testing.T, createdEvent *remote.Event, err error) {
				require.NoError(t, err)
				require.NotNil(t, createdEvent)
				require.Equal(t, &remote.Event{}, createdEvent)
			},
		},
		{
			name: "error creating event",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUserID: "testMMUserID",
			},
			event:             &remote.Event{Subject: "Test Event"},
			mattermostUserIDs: []string{"testMMUserID1"},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID1").Return(nil, errors.New("not found")).Times(1)
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockPoster.EXPECT().DM("testMMUserID1", gomock.Any(), "testDisplayName", "testDisplayName", "testCommandTrigger").Return("", fmt.Errorf("error creating DM")).Times(1).Return("", nil)
				mockClient.EXPECT().CreateEvent("testRemoteUserID", &remote.Event{Subject: "Test Event"}).Return(nil, fmt.Errorf("error creating event")).Times(1)
			},
			assertions: func(t *testing.T, createdEvent *remote.Event, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error creating event")
			},
		},
		{
			name: "successful event creation",
			user: &User{
				User:             &store.User{Remote: &remote.User{ID: "testRemoteUserID"}},
				MattermostUserID: "testMMUserID",
			},
			event: &remote.Event{
				Subject:   "Test Event",
				Location:  &remote.Location{DisplayName: "Test Location"},
				Start:     &remote.DateTime{DateTime: "2024-10-01T09:00:00", TimeZone: "UTC"},
				End:       &remote.DateTime{DateTime: "2024-10-01T10:00:00", TimeZone: "UTC"},
				Attendees: []*remote.Attendee{{EmailAddress: &remote.EmailAddress{Address: "attendee1@example.com"}}},
			},
			mattermostUserIDs: []string{"testMMUserID1"},
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID1").Return(nil, errors.New("not found")).Times(1)
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockPoster.EXPECT().DM("testMMUserID1", gomock.Any(), "testDisplayName", "testDisplayName", "testCommandTrigger").Return("", fmt.Errorf("error creating DM")).Times(1).Return("", nil)
				mockClient.EXPECT().CreateEvent("testRemoteUserID", &remote.Event{
					Subject:   "Test Event",
					Location:  &remote.Location{DisplayName: "Test Location"},
					Start:     &remote.DateTime{DateTime: "2024-10-01T09:00:00", TimeZone: "UTC"},
					End:       &remote.DateTime{DateTime: "2024-10-01T10:00:00", TimeZone: "UTC"},
					Attendees: []*remote.Attendee{{EmailAddress: &remote.EmailAddress{Address: "attendee1@example.com"}}},
				}).Return(&remote.Event{Subject: "Created Test Event", ID: "123"}, nil).Times(1)
			},
			assertions: func(t *testing.T, createdEvent *remote.Event, err error) {
				require.NoError(t, err)
				require.NotNil(t, createdEvent)
				require.Equal(t, &remote.Event{Subject: "Created Test Event", ID: "123"}, createdEvent)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			createdEvent, err := mscalendar.CreateEvent(tt.user, tt.event, tt.mattermostUserIDs)
			tt.assertions(t, createdEvent, err)
		})
	}
}

func TestDeleteCalendar(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	user := &User{
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name       string
		calendarID string
		setupMock  func()
		assertions func(t *testing.T, err error)
	}{
		{
			name:       "error filtering with client",
			calendarID: "testCalendarID",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
			},
		},
		{
			name:       "error deleting calendar",
			calendarID: "testCalendarID",
			setupMock: func() {
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().DeleteCalendar(user.User.Remote.ID, "testCalendarID").Return(errors.New("deletion error")).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "deletion error")
			},
		},
		{
			name:       "successful calendar deletion",
			calendarID: "testCalendarID",
			setupMock: func() {
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockClient.EXPECT().DeleteCalendar(user.User.Remote.ID, "testCalendarID").Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.DeleteCalendar(user, tt.calendarID)

			tt.assertions(t, err)
		})
	}
}

func TestFindMeetingTimes(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	user := &User{
		MattermostUserID: "testMMUserID",
	}

	meetingParams := &remote.FindMeetingTimesParameters{}

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error, results *remote.MeetingTimeSuggestionResults)
	}{
		{
			name: "error filtering with client",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, err error, results *remote.MeetingTimeSuggestionResults) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
				require.Nil(t, results)
			},
		},
		{
			name: "error finding meeting times",
			setupMock: func() {
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().FindMeetingTimes(user.User.Remote.ID, meetingParams).Return(nil, errors.New("finding times error")).Times(1)
			},
			assertions: func(t *testing.T, err error, results *remote.MeetingTimeSuggestionResults) {
				require.Error(t, err)
				require.EqualError(t, err, "finding times error")
				require.Nil(t, results)
			},
		},
		{
			name: "successful meeting time retrieval",
			setupMock: func() {
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().FindMeetingTimes(user.User.Remote.ID, meetingParams).Return(&remote.MeetingTimeSuggestionResults{}, nil).Times(1)
			},
			assertions: func(t *testing.T, err error, results *remote.MeetingTimeSuggestionResults) {
				require.NoError(t, err)
				require.NotNil(t, results)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			results, err := mscalendar.FindMeetingTimes(user, meetingParams)

			tt.assertions(t, err, results)
		})
	}
}

func TestGetCalendars(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	user := &User{
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error, calendars []*remote.Calendar)
	}{
		{
			name: "error filtering with client",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error loading the user")).Times(1)
			},
			assertions: func(t *testing.T, err error, calendars []*remote.Calendar) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error loading the user")
				require.Nil(t, calendars)
			},
		},
		{
			name: "error getting calendars",
			setupMock: func() {
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().GetCalendars(user.User.Remote.ID).Return(nil, errors.New("getting calendars error")).Times(1)
			},
			assertions: func(t *testing.T, err error, calendars []*remote.Calendar) {
				require.Error(t, err)
				require.EqualError(t, err, "getting calendars error")
				require.Nil(t, calendars)
			},
		},
		{
			name: "successful calendars retrieval",
			setupMock: func() {
				user.User = &store.User{Remote: &remote.User{ID: "testRemoteUserID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().GetCalendars(user.User.Remote.ID).Return([]*remote.Calendar{{ID: "calendar1"}}, nil).Times(1)
			},
			assertions: func(t *testing.T, err error, calendars []*remote.Calendar) {
				require.NoError(t, err)
				require.Equal(t, []*remote.Calendar{{ID: "calendar1"}}, calendars)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			calendars, err := mscalendar.GetCalendars(user)

			tt.assertions(t, err, calendars)
		})
	}
}

// revive:disable:unexported-return
func MockSetup(t *testing.T) (*mscalendar, *mock_store.MockStore, *mock_bot.MockPoster, *mock_remote.MockRemote, *mock_plugin_api.MockPluginAPI, *mock_remote.MockClient, *mock_bot.MockLogger) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	mockPoster := mock_bot.NewMockPoster(ctrl)
	mockRemote := mock_remote.NewMockRemote(ctrl)
	mockPluginAPI := mock_plugin_api.NewMockPluginAPI(ctrl)
	mockClient := mock_remote.NewMockClient(ctrl)
	mockLogger := mock_bot.NewMockLogger(ctrl)

	env := Env{
		Dependencies: &Dependencies{
			Store:     mockStore,
			Poster:    mockPoster,
			Remote:    mockRemote,
			PluginAPI: mockPluginAPI,
			Logger:    mockLogger,
		},
	}

	mscalendar := &mscalendar{
		Env:    env,
		client: mockClient,
	}

	mscalendar.Config = &config.Config{
		Provider: config.ProviderConfig{
			DisplayName:    "testDisplayName",
			CommandTrigger: "testCommandTrigger",
		},
	}

	return mscalendar, mockStore, mockPoster, mockRemote, mockPluginAPI, mockClient, mockLogger
}
