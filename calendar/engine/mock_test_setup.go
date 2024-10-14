package engine

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/config"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/engine/mock_plugin_api"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote/mock_remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store/mock_store"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/utils/bot/mock_bot"
)

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
		PluginVersion: "1.0.0",
	}

	return mscalendar, mockStore, mockPoster, mockRemote, mockPluginAPI, mockClient, mockLogger
}

func GetMockSubscription() *store.Subscription {
	return &store.Subscription{
		Remote:              &remote.Subscription{},
		MattermostCreatorID: "testActingUserID",
		PluginVersion:       "1.0.0",
	}
}

func GetMockUser() *User {
	return &User{
		MattermostUserID: "testMMUserID",
	}
}
