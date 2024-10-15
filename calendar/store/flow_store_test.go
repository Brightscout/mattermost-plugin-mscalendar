package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestSetProperty(t *testing.T) {
	mockAPI, store, _, _, mockTracker := MockStoreSetup(t)
	mockUser := User{
		MattermostUserID: "mockMMUserID",
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		Settings: Settings{
			DailySummary: &DailySummaryUserSettings{
				PostTime: "10:00AM",
			},
			EventSubscriptionID:               "mockEventSubscriptionID",
			UpdateStatusFromOptions:           "available",
			GetConfirmation:                   true,
			ReceiveReminders:                  true,
			SetCustomStatus:                   false,
			UpdateStatus:                      false,
			ReceiveNotificationsDuringMeeting: true,
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name         string
		propertyName string
		value        interface{}
		setup        func()
		assertions   func(*testing.T, error)
	}{
		{
			name:         "Error loading user",
			propertyName: UpdateStatusPropertyName,
			value:        "online",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
			},
		},
		{
			name:         "Set UpdateStatusPropertyName successfully",
			propertyName: UpdateStatusPropertyName,
			value:        "online",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "online", "flow").Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Set GetConfirmationPropertyName successfully",
			propertyName: GetConfirmationPropertyName,
			value:        true,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Set SetCustomStatusPropertyName successfully",
			propertyName: SetCustomStatusPropertyName,
			value:        false,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Set ReceiveUpcomingEventReminderName successfully",
			propertyName: ReceiveUpcomingEventReminderName,
			value:        false,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Invalid property name",
			propertyName: "mockPropertyName",
			value:        false,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "property mockPropertyName not found")
			},
		},
		{
			name:         "Error storing user",
			propertyName: SetCustomStatusPropertyName,
			value:        true,
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.SetProperty("mockUserID", tt.propertyName, tt.value)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestSetPostID(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			PostIDs: make(map[string]string),
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name         string
		userID       string
		propertyName string
		postID       string
		setup        func()
		assertions   func(*testing.T, error)
	}{
		{
			name:         "Error loading user",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			postID:       "mockPostID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
			},
		},
		{
			name:         "Set PostID successfully",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			postID:       "mockPostID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Set PostID for user with nil PostIDs map",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			postID:       "mockPostID",
			setup: func() {
				mockUser.WelcomeFlowStatus.PostIDs = nil
				mockUserJSON, _ = json.Marshal(mockUser)
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Error storing user",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			postID:       "mockPostID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.SetPostID(tt.userID, tt.propertyName, tt.postID)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestGetPostID(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			PostIDs: map[string]string{
				"welcomePost": "mockPostID",
			},
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name         string
		userID       string
		propertyName string
		setup        func()
		assertions   func(*testing.T, string, error)
	}{
		{
			name:         "Error loading user",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, postID string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
				require.Empty(t, postID)
			},
		},
		{
			name:         "PostID retrieved successfully",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, postID string, err error) {
				require.NoError(t, err)
				require.Equal(t, "mockPostID", postID)
			},
		},
		{
			name:         "PostID does not exist",
			userID:       "mockUserID",
			propertyName: "nonExistentPost",
			setup: func() {
				mockUser.WelcomeFlowStatus.PostIDs = map[string]string{}
				mockUserJSON, _ = json.Marshal(mockUser)
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, postID string, err error) {
				require.NoError(t, err)
				require.Empty(t, postID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			postID, err := store.GetPostID(tt.userID, tt.propertyName)

			tt.assertions(t, postID, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestRemovePostID(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			PostIDs: map[string]string{
				"welcomePost": "mockPostID",
			},
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name         string
		userID       string
		propertyName string
		setup        func()
		assertions   func(*testing.T, error)
	}{
		{
			name:         "Error loading user",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
			},
		},
		{
			name:         "Remove PostID successfully",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:         "Error storing user",
			userID:       "mockUserID",
			propertyName: "welcomePost",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.RemovePostID(tt.userID, tt.propertyName)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestGetCurrentStep(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			Step: 3,
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		userID     string
		setup      func()
		assertions func(*testing.T, int, error)
	}{
		{
			name:   "Error loading user",
			userID: "mockUserID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, step int, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
				require.Equal(t, 0, step)
			},
		},
		{
			name:   "Get current step successfully",
			userID: "mockUserID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, step int, err error) {
				require.NoError(t, err)
				require.Equal(t, 3, step)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			step, err := store.GetCurrentStep(tt.userID)

			tt.assertions(t, step, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestSetCurrentStep(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			Step: 1,
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		userID     string
		step       int
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name:   "Error loading user",
			userID: "mockUserID",
			step:   2,
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
			},
		},
		{
			name:   "Error storing user",
			userID: "mockUserID",
			step:   2,
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
		{
			name:   "Set current step successfully",
			userID: "mockUserID",
			step:   2,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.SetCurrentStep(tt.userID, tt.step)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteCurrentStep(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)
	mockUser := User{
		Remote:           &remote.User{ID: "mockRemoteUserID"},
		MattermostUserID: "mockMMUserID",
		WelcomeFlowStatus: WelcomeFlowStatus{
			Step: 1,
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		userID     string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name:   "Error loading user",
			userID: "mockUserID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Error loading user")
			},
		},
		{
			name:   "Error storing user",
			userID: "mockUserID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
		{
			name:   "Delete current step successfully",
			userID: "mockUserID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := store.DeleteCurrentStep(tt.userID)
			tt.assertions(t, err)
			mockAPI.AssertExpectations(t)
		})
	}
}
