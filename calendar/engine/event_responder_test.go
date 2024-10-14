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
		name      string
		setupMock func()
		assertion func(err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error accepting event",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().AcceptEvent("testRemoteID", "mockEventID").Return(errors.New("unable to accept event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to accept event")
			},
		},
		{
			name: "successful event acceptance",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().AcceptEvent("testRemoteID", "mockEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.AcceptEvent(user, "mockEventID")

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
		name      string
		setupMock func()
		assertion func(err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error declining event",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(errors.New("unable to decline event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to decline event")
			},
		},
		{
			name: "successful event decline",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.DeclineEvent(user, "testEventID")

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
		name      string
		setupMock func()
		assertion func(err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error tentatively accepting event",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to tentatively accept event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to tentatively accept event")
			},
		},
		{
			name: "successful tentative event acceptance",
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID")
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.TentativelyAcceptEvent(user, "testEventID")

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
		name      string
		response  string
		setupMock func()
		assertion func(err error)
	}{
		{
			name:      "invalid response error",
			response:  OptionNotResponded,
			setupMock: func() {},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "not responded is not a valid response")
			},
		},
		{
			name:     "invalid response string",
			response: "InvalidResponse",
			setupMock: func() {
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(&model.User{Id: "testMMUserID"}, nil)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "InvalidResponse is not a valid response")
			},
		},
		{
			name:     "error filtering user",
			response: OptionYes,
			setupMock: func() {
				user.User = nil
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name:     "success accepting event",
			response: OptionYes,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error accepting event",
			response: OptionYes,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().AcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to accept event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to accept event")
			},
		},
		{
			name:     "success declining event",
			response: OptionNo,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error declining event",
			response: OptionNo,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().DeclineEvent("testRemoteID", "testEventID").Return(errors.New("unable to decline event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to decline event")
			},
		},
		{
			name:     "success tentatively accepting event",
			response: OptionMaybe,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(nil).Times(1)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "error tentatively accepting event",
			response: OptionMaybe,
			setupMock: func() {
				user.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().TentativelyAcceptEvent("testRemoteID", "testEventID").Return(errors.New("unable to tentatively accept event")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "unable to tentatively accept event")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.RespondToEvent(user, "testEventID", tt.response)

			tt.assertion(err)
		})
	}
}
