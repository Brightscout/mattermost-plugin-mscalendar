package engine

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestAcceptEvent(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	user := &User{
		User:             &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}},
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name                 string
		eventID              string
		setupMock            func()
		expectedErrorMessage string
		assertion            func(err error)
	}{
		{
			name:    "error filtering with user",
			eventID: "testEventID",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			expectedErrorMessage: "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name:    "error accepting event",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to accept event")).Times(1)
			},
			expectedErrorMessage: "unable to accept event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to accept event")
			},
		},
		{
			name:    "successful event acceptance",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.AcceptEvent(user, tt.eventID)

			tt.assertion(err)
		})
	}
}

func TestDeclineEvent(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	user := &User{
		User:             &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}},
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name                 string
		eventID              string
		setupMock            func()
		expectedErrorMessage string
		assertion            func(err error)
	}{
		{
			name:    "error filtering with user",
			eventID: "testEventID",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			expectedErrorMessage: "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name:    "error declining event",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(errors.New("unable to decline event")).Times(1)
			},
			expectedErrorMessage: "unable to decline event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to decline event")
			},
		},
		{
			name:    "successful event decline",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.DeclineEvent(user, tt.eventID)

			tt.assertion(err)
		})
	}
}

func TestTentativelyAcceptEvent(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	user := &User{
		User:             &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}},
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name                 string
		eventID              string
		setupMock            func()
		expectedErrorMessage string
		assertion            func(err error)
	}{
		{
			name:    "error filtering with user",
			eventID: "testEventID",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			expectedErrorMessage: "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name:    "error tentatively accepting event",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to tentatively accept event")).Times(1)
			},
			expectedErrorMessage: "unable to tentatively accept event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to tentatively accept event")
			},
		},
		{
			name:    "successful tentative event acceptance",
			eventID: "testEventID",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.TentativelyAcceptEvent(user, tt.eventID)

			tt.assertion(err)
		})
	}
}

func TestRespondToEvent(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)

	user := &User{
		User:             &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}},
		MattermostUserID: "testMMUserID",
	}

	tests := []struct {
		name                 string
		eventID              string
		response             string
		setupMock            func()
		expectedErrorMessage string
		assertion            func(err error)
	}{
		{
			name:                 "error - not responded is an invalid response",
			eventID:              "testEventID",
			response:             OptionNotResponded,
			setupMock:            func() {},
			expectedErrorMessage: "not responded is not a valid response",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "not responded is not a valid response")
			},
		},
		{
			name:     "error - invalid response string",
			eventID:  "testEventID",
			response: "InvalidResponse",
			setupMock: func() {
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(&model.User{Id: "testMMUserID"}, nil)
			},
			expectedErrorMessage: "InvalidResponse is not a valid response",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "InvalidResponse is not a valid response")
			},
		},
		{
			name:     "error - filtering fails",
			eventID:  "testEventID",
			response: OptionYes,
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			expectedErrorMessage: "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name:     "success - accepted event",
			eventID:  "testEventID",
			response: OptionYes,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error - accept event fails",
			eventID:  "testEventID",
			response: OptionYes,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to accept event")).Times(1)
			},
			expectedErrorMessage: "unable to accept event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to accept event")
			},
		},
		{
			name:     "success - declined event",
			eventID:  "testEventID",
			response: OptionNo,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error - decline event fails",
			eventID:  "testEventID",
			response: OptionNo,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(errors.New("unable to decline event")).Times(1)
			},
			expectedErrorMessage: "unable to decline event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to decline event")
			},
		},
		{
			name:     "success - tentatively accepted event",
			eventID:  "testEventID",
			response: OptionMaybe,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			expectedErrorMessage: "",
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error - tentatively accept event fails",
			eventID:  "testEventID",
			response: OptionMaybe,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to tentatively accept event")).Times(1)
			},
			expectedErrorMessage: "unable to tentatively accept event",
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to tentatively accept event")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.RespondToEvent(user, tt.eventID, tt.response)

			tt.assertion(err)
		})
	}
}
