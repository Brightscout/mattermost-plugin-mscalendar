package store

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/testutil"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/tracker/mock_tracker"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/utils/bot/mock_bot"
)

const (
	MockEventSubscriptionID = "mockEventSubscriptionID"
	MockSubscriptionID      = "mockSubscriptionID"
	MockRemoteUserID        = "mockRemoteUserID"
	MockRemoteID            = "mockRemoteID"
	MockCreatorID           = "mockCreatorID"
	MockMMUserID            = "mockMMUserID"
)

func GetMockSetup(t *testing.T) (*testutil.MockPluginAPI, Store, *mock_bot.MockLogger, *mock_bot.MockLogger, *mock_tracker.MockTracker) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mock_bot.NewMockLogger(ctrl)
	mockLoggerWith := mock_bot.NewMockLogger(ctrl)
	mockTracker := mock_tracker.NewMockTracker(ctrl)
	mockAPI := &testutil.MockPluginAPI{}
	store := NewPluginStore(mockAPI, mockLogger, mockTracker, false, nil)

	return mockAPI, store, mockLogger, mockLoggerWith, mockTracker
}

func GetMockUser() *User {
	return &User{
		MattermostUserID: MockMMUserID,
		Settings: Settings{
			EventSubscriptionID: MockEventSubscriptionID,
		},
		Remote: &remote.User{
			ID: MockRemoteUserID,
		},
	}
}

func GetMockSubscription() *Subscription {
	return &Subscription{
		Remote: &remote.Subscription{
			ID:        MockSubscriptionID,
			CreatorID: MockCreatorID,
		},
	}
}
