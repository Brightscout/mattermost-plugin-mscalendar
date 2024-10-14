package engine

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/config"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/engine/mock_plugin_api"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/engine/mock_welcomer"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote/mock_remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store/mock_store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/utils/bot/mock_bot"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestExpandUser(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, _, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error)
	}{
		{
			name: "error expanding remote user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error expanding Mattermost user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(&store.User{}, nil).Times(1)
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(nil, errors.New("some error occurred while getting Mattermost user"))
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "some error occurred while getting Mattermost user")
			},
		},
		{
			name: "success expanding user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(&store.User{}, nil).Times(1)
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(&model.User{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.ExpandUser(mockUser)

			tt.assertions(t, err)
		})
	}
}

func TestExpandRemoteUser(t *testing.T) {
	mscalendar, mockStore, _, _, _, _, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error)
	}{
		{
			name: "error loading remote user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "success expanding remote user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(&store.User{}, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.ExpandRemoteUser(mockUser)

			tt.assertions(t, err)
		})
	}
}

func TestExpandMattermostUser(t *testing.T) {
	mscalendar, _, _, _, mockPluginAPI, _, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error)
	}{
		{
			name: "error expanding Mattermost user",
			setupMock: func() {
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(nil, errors.New("some error occurred while getting Mattermost user"))
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "some error occurred while getting Mattermost user")
			},
		},
		{
			name: "success expanding Mattermost user",
			setupMock: func() {
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(&model.User{}, nil)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.ExpandMattermostUser(mockUser)

			tt.assertions(t, err)
		})
	}
}

func TestGetTimezone(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(t *testing.T, err error)
	}{
		{
			name: "error loading remote user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name: "error getting mailbox setting",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().GetMailboxSettings("testRemoteID").Return(nil, errors.New("error occurred getting mailbox settings"))
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error occurred getting mailbox settings")
			},
		},
		{
			name: "success getting mailbox setting",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockClient.EXPECT().GetMailboxSettings("testRemoteID").Return(&remote.MailboxSettings{TimeZone: "mockTimeZone"}, nil)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			_, err := mscalendar.GetTimezone(mockUser)

			tt.assertions(t, err)
		})
	}
}

func TestUser_String(t *testing.T) {
	tests := []struct {
		name       string
		user       *User
		assertions func(t *testing.T, actualString string)
	}{
		{
			name: "User with Mattermost user object",
			user: &User{
				MattermostUserID: "mockMMUserID",
				MattermostUser: &model.User{
					Username: "mockMMUsername",
				},
			},
			assertions: func(t *testing.T, actualString string) {
				require.Equal(t, "@mockMMUsername", actualString)
			},
		},
		{
			name: "User without Mattermost user object",
			user: &User{
				MattermostUserID: "mockMMUserID",
			},
			assertions: func(t *testing.T, actualString string) {
				require.Equal(t, "mockMMUserID", actualString)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualString := tt.user.String()

			tt.assertions(t, actualString)
		})
	}
}

func TestUser_Markdown(t *testing.T) {
	tests := []struct {
		name       string
		user       *User
		assertions func(t *testing.T, actualOutput string)
	}{
		{
			name: "User with Mattermost user object",
			user: &User{
				MattermostUserID: "mockMMUserID",
				MattermostUser: &model.User{
					Username: "mockMMUsername",
				},
			},
			assertions: func(t *testing.T, actualOutput string) {
				require.Equal(t, "@mockMMUsername", actualOutput)
			},
		},
		{
			name: "User without Mattermost user object",
			user: &User{
				MattermostUserID: "mockMMUserID",
			},
			assertions: func(t *testing.T, actualOutput string) {
				require.Equal(t, "UserID: `mockMMUserID`", actualOutput)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOutput := tt.user.Markdown()

			tt.assertions(t, actualOutput)
		})
	}
}

func TestDisconnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	mockPoster := mock_bot.NewMockPoster(ctrl)
	mockRemote := mock_remote.NewMockRemote(ctrl)
	mockPluginAPI := mock_plugin_api.NewMockPluginAPI(ctrl)
	mockClient := mock_remote.NewMockClient(ctrl)
	mockLogger := mock_bot.NewMockLogger(ctrl)
	mockLoggerWith := mock_bot.NewMockLogger(ctrl)
	mockWelcomer := mock_welcomer.NewMockWelcomer(ctrl)

	env := Env{
		Dependencies: &Dependencies{
			Store:     mockStore,
			Poster:    mockPoster,
			Remote:    mockRemote,
			PluginAPI: mockPluginAPI,
			Logger:    mockLogger,
			Welcomer:  mockWelcomer,
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

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(err error)
	}{
		{
			name: "error filtering user",
			setupMock: func() {
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockRemoteMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name: "error loading user",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(nil, errors.New("error loading user")).Times(1)
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error loading user")
			},
		},
		{
			name: "error deleting linked channels from events",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{ChannelEvents: store.ChannelEventLink{"mockEventID": "mockChannelID"}, MattermostDisplayName: "mockMMUserDisplayName"}, nil).Times(1)
				mockStore.EXPECT().DeleteLinkedChannelFromEvent("mockEventID", "mockChannelID").Return(errors.New("some error occurred deleting linked channel"))
				mockStore.EXPECT().StoreUser(gomock.Any()).Return(errors.New("some error occurred storing user"))
				mockLogger.EXPECT().With(gomock.Any()).Return(mockLoggerWith).Times(1)
				mockLoggerWith.EXPECT().Errorf("error storing user after failing deleting linked channels from store").Times(1)
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error deleting linked channels from events")
			},
		},
		{
			name: "error loading subscription",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Settings: store.Settings{EventSubscriptionID: "mockEventSubscriptionID"}}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("mockEventSubscriptionID").Return(nil, errors.New("internal error"))
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error loading subscription: internal error")
			},
		},
		{
			name: "failed to delete event subscription",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Settings: store.Settings{EventSubscriptionID: "mockEventSubscriptionID"}}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("mockEventSubscriptionID").Return(nil, nil)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "mockEventSubscriptionID").Return(errors.New("internal server error"))
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed to delete subscription mockEventSubscriptionID: internal server error")
			},
		},
		{
			name: "error deleting user",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Settings: store.Settings{EventSubscriptionID: "mockEventSubscriptionID"}}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("mockEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "mockEventSubscriptionID").Return(nil)
				mockClient.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
				mockStore.EXPECT().DeleteUser("mockMMUserID").Return(errors.New("error deleting user"))
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error deleting user")
			},
		},
		{
			name: "error deleting user from index",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Settings: store.Settings{EventSubscriptionID: "mockEventSubscriptionID"}}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("mockEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "mockEventSubscriptionID").Return(nil)
				mockClient.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
				mockStore.EXPECT().DeleteUser("mockMMUserID").Return(nil)
				mockStore.EXPECT().DeleteUserFromIndex("mockMMUserID").Return(errors.New("error deleting user from index"))
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error deleting user from index")
			},
		},
		{
			name: "user disconnected successfully",
			setupMock: func() {
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{MattermostUserID: "mockRemoteMMUserID"}
				mockWelcomer.EXPECT().AfterDisconnect("mockMMUserID").Return(nil)
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Settings: store.Settings{EventSubscriptionID: "mockEventSubscriptionID"}}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("mockEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "mockEventSubscriptionID").Return(nil)
				mockClient.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
				mockStore.EXPECT().DeleteUser("mockMMUserID").Return(nil)
				mockStore.EXPECT().DeleteUserFromIndex("mockMMUserID").Return(nil)
			},
			assertions: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.DisconnectUser("mockMMUserID")

			tt.assertions(err)
		})
	}
}

func TestGetRemoteUser(t *testing.T) {
	mscalendar, mockStore, _, _, _, _, _ := MockSetup(t)

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(remoteUser *remote.User, err error)
	}{
		{
			name: "LoadUser returns an error",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(nil, errors.New("failed to load user")).Times(1)
			},
			assertions: func(remoteUser *remote.User, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed to load user")
				require.Nil(t, remoteUser)
			},
		},
		{
			name: "Successfully get remote user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("mockMMUserID").Return(&store.User{Remote: &remote.User{ID: "mockRemoteUserID"}}, nil).Times(1)
			},
			assertions: func(remoteUser *remote.User, err error) {
				require.NoError(t, err)
				require.Equal(t, &remote.User{ID: "mockRemoteUserID"}, remoteUser)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			remoteUser, err := mscalendar.GetRemoteUser("mockMMUserID")

			tt.assertions(remoteUser, err)
		})
	}
}

func TestIsAuthorizedAdmin(t *testing.T) {
	mscalendar, _, _, _, mockPluginAPI, _, _ := MockSetup(t)

	tests := []struct {
		name             string
		mattermostUserID string
		setupMock        func()
		assertions       func(result bool, err error)
	}{
		{
			name:             "User is in AdminUserIDs",
			mattermostUserID: "mockAdminID1",
			setupMock: func() {
				mscalendar.AdminUserIDs = "mockAdminID1,mockAdminID2"
			},
			assertions: func(result bool, err error) {
				require.NoError(t, err)
				require.Equal(t, true, result)
			},
		},
		{
			name:             "error checking system admin",
			mattermostUserID: "mockMMUserID",
			setupMock: func() {
				mscalendar.AdminUserIDs = "mockAdminID1,mockAdminID2"
				mockPluginAPI.EXPECT().IsSysAdmin("mockMMUserID").Return(false, errors.New("error occurred checking system admin")).Times(1)
			},
			assertions: func(result bool, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error occurred checking system admin")
			},
		},
		{
			name:             "User is not in AdminUserIDs and is not a system admin",
			mattermostUserID: "mockMMUserID",
			setupMock: func() {
				mscalendar.AdminUserIDs = "mockAdminID1,mockAdminID2"
				mockPluginAPI.EXPECT().IsSysAdmin("mockMMUserID").Return(false, nil).Times(1)
			},
			assertions: func(result bool, err error) {
				require.NoError(t, err)
				require.Equal(t, false, result)
			},
		},
		{
			name:             "User is not in AdminUserIDs but is a system admin",
			mattermostUserID: "mockMMUserID",
			setupMock: func() {
				mscalendar.AdminUserIDs = "mockAdminID1,mockAdminID2"
				mockPluginAPI.EXPECT().IsSysAdmin("mockMMUserID").Return(true, nil).Times(1)
			},
			assertions: func(result bool, err error) {
				require.NoError(t, err)
				require.Equal(t, true, result)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := mscalendar.IsAuthorizedAdmin(tt.mattermostUserID)

			tt.assertions(result, err)
		})
	}
}

func TestGetUserSettings(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, _, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setupMock  func()
		assertions func(result *store.Settings, err error)
	}{
		{
			name: "error filtering the user",
			setupMock: func() {
				mockStore.EXPECT().LoadUser("testMMUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertions: func(result *store.Settings, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "Successfully get user settings",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{GetConfirmation: false}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testMMUserID").Return(&model.User{}, nil)
			},
			assertions: func(result *store.Settings, err error) {
				require.NoError(t, err)
				require.Equal(t, &store.Settings{GetConfirmation: false}, result)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := mscalendar.GetUserSettings(mockUser)

			tt.assertions(result, err)
		})
	}
}
