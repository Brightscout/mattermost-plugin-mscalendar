package store

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/testutil"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/utils/bot/mock_bot"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestLoadSubscription(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, *Subscription, error)
	}{
		{
			name: "Error loading subscription",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Subscription not found"}).Times(1)
			},
			assertions: func(t *testing.T, sub *Subscription, err error) {
				require.Error(t, err)
				require.Nil(t, sub)
				require.EqualError(t, err, "failed plugin KVGet: Subscription not found")
			},
		},
		{
			name: "Successful Load",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"PluginVersion":"1.0","Remote":{"ID":"mockRemoteUserID","CreatorID":"mockCreatorID"}}`), nil).Times(1)
			},
			assertions: func(t *testing.T, sub *Subscription, err error) {
				require.NoError(t, err)
				require.NotNil(t, sub)
				require.Equal(t, "1.0", sub.PluginVersion)
				require.Equal(t, MockRemoteUserID, sub.Remote.ID)
				require.Equal(t, MockCreatorID, sub.Remote.CreatorID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI, store, _, _, _ := GetMockSetup(t)
			tt.setup(mockAPI)

			sub, err := store.LoadSubscription(MockSubscriptionID)

			tt.assertions(t, sub, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserSubscription(t *testing.T) {
	mockUser := GetMockUser()
	mockSubscription := GetMockSubscription()

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI, *mock_bot.MockLogger, *mock_bot.MockLogger)
		assertions func(*testing.T, error)
	}{
		{
			name:  "User does not match subscription creator",
			setup: func(_ *testutil.MockPluginAPI, _ *mock_bot.MockLogger, _ *mock_bot.MockLogger) {},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, `user "mockRemoteID" does not match the subscription creator "mockCreatorID"`)
			},
		},
		{
			name: "Error storing subscription",
			setup: func(mockAPI *testutil.MockPluginAPI, _ *mock_bot.MockLogger, _ *mock_bot.MockLogger) {
				mockSubscription.Remote.CreatorID = mockUser.Remote.ID
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store subscription"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to store subscription")
			},
		},
		{
			name: "Error storing user settings",
			setup: func(mockAPI *testutil.MockPluginAPI, _ *mock_bot.MockLogger, _ *mock_bot.MockLogger) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store user settings"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to store user settings")
			},
		},
		{
			name: "Successful Store",
			setup: func(mockAPI *testutil.MockPluginAPI, mockLogger *mock_bot.MockLogger, mockLoggerWith *mock_bot.MockLogger) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockLogger.EXPECT().With(gomock.Any()).Return(mockLoggerWith).Times(1)
				mockLoggerWith.EXPECT().Debugf("store: stored mattermost user subscription.").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI, store, mockLogger, mockLoggerWith, _ := GetMockSetup(t)
			tt.setup(mockAPI, mockLogger, mockLoggerWith)

			err := store.StoreUserSubscription(mockUser, mockSubscription)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteUserSubscription(t *testing.T) {
	mockUser := GetMockUser()

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI, *mock_bot.MockLogger, *mock_bot.MockLogger)
		assertions func(*testing.T, error)
	}{
		{
			name: "Error deleting subscription",
			setup: func(mockAPI *testutil.MockPluginAPI, _ *mock_bot.MockLogger, _ *mock_bot.MockLogger) {
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(&model.AppError{Message: "Failed to delete subscription"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to delete subscription")
			},
		},
		{
			name: "Error updating user settings",
			setup: func(mockAPI *testutil.MockPluginAPI, _ *mock_bot.MockLogger, _ *mock_bot.MockLogger) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to update user settings"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to update user settings")
			},
		},
		{
			name: "Successful Delete",
			setup: func(mockAPI *testutil.MockPluginAPI, mockLogger *mock_bot.MockLogger, mockLoggerWith *mock_bot.MockLogger) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockLogger.EXPECT().With(gomock.Any()).Return(mockLoggerWith).Times(1)
				mockLoggerWith.EXPECT().Debugf("store: deleted mattermost user subscription.").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI, store, mockLogger, mockLoggerWith, _ := GetMockSetup(t)
			tt.setup(mockAPI, mockLogger, mockLoggerWith)

			err := store.DeleteUserSubscription(mockUser, MockSubscriptionID)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}
