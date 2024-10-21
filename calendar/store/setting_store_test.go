package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
)

func TestSetSetting(t *testing.T) {
	mockAPI, store, _, _, mockTracker := GetMockSetup(t)
	mockUser := User{
		MattermostUserID: "mockMattermostUserID",
		Remote: &remote.User{
			ID: "mockRemoteUserID",
		},
	}
	mockUserJSON, err := json.Marshal(mockUser)
	require.NoError(t, err)

	tests := []struct {
		name       string
		settingID  string
		value      interface{}
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name:      "Error loading user",
			settingID: "mockSettingID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading user"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "failed plugin KVGet: Error loading user")
			},
		},
		{
			name:      "error setting UpdateStatusFromOptionsSetting",
			settingID: UpdateStatusFromOptionsSettingID,
			value:     1,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "cannot read value 1 for setting update_status_from_options (expecting string)")
			},
		},
		{
			name:      "Set UpdateStatusFromOptionsSetting",
			settingID: UpdateStatusFromOptionsSettingID,
			value:     "available",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "error setting GetConfirmationSettingID",
			settingID: GetConfirmationSettingID,
			value:     1,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "cannot read value 1 for setting get_confirmation (expecting bool)")
			},
		},
		{
			name:      "Set GetConfirmationSettingID",
			settingID: GetConfirmationSettingID,
			value:     true,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "error setting SetCustomStatusSettingID",
			settingID: SetCustomStatusSettingID,
			value:     1,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "cannot read value 1 for setting set_custom_status (expecting bool)")
			},
		},
		{
			name:      "Set SetCustomStatusSettingID",
			settingID: SetCustomStatusSettingID,
			value:     true,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "error setting ReceiveRemindersSettingID",
			settingID: ReceiveRemindersSettingID,
			value:     1,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "cannot read value 1 for setting get_reminders (expecting bool)")
			},
		},
		{
			name:      "Set ReceiveRemindersSettingID",
			settingID: ReceiveRemindersSettingID,
			value:     true,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "Set DailySummarySettingID",
			settingID: DailySummarySettingID,
			value:     "mockDailySummarySetting",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "invalid setting ID",
			settingID: "invalidSettingID",
			value:     "mockDailySummarySetting",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "setting invalidSettingID not found")
			},
		},
		{
			name:      "Error storing updated user",
			settingID: UpdateStatusFromOptionsSettingID,
			value:     "available",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Error storing user"}).Times(1)
				mockTracker.EXPECT().TrackAutomaticStatusUpdate("mockUserID", "available", "settings").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Error storing user")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.SetSetting("mockUserID", tt.settingID, tt.value)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestGetSetting(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)
	mockUser := User{
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
		name       string
		settingID  string
		setup      func()
		assertions func(*testing.T, interface{}, error)
	}{
		{
			name:      "Error loading settings",
			settingID: "mockSettingID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading settings"}).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.Error(t, err)
				require.Nil(t, setting)
				require.EqualError(t, err, "failed plugin KVGet: Error loading settings")
			},
		},
		{
			name:      "Get UpdateStatusFromOptionsSetting",
			settingID: UpdateStatusFromOptionsSettingID,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, "available", setting)
			},
		},
		{
			name:      "Get GetConfirmationSetting",
			settingID: GetConfirmationSettingID,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, true, setting)
			},
		},
		{
			name:      "Get SetCustomStatusSetting",
			settingID: SetCustomStatusSettingID,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, false, setting)
			},
		},
		{
			name:      "Get ReceiveRemindersSetting",
			settingID: ReceiveRemindersSettingID,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, true, setting)
			},
		},
		{
			name:      "Get DailySummary",
			settingID: DailySummarySettingID,
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, &DailySummaryUserSettings{PostTime: "10:00AM"}, setting)
			},
		},
		{
			name:      "invalid settingID",
			settingID: "invalidSettingID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(mockUserJSON, nil).Times(1)
			},
			assertions: func(t *testing.T, setting interface{}, err error) {
				require.EqualError(t, err, "setting invalidSettingID not found")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			setting, err := store.GetSetting("mockUserID", tt.settingID)

			tt.assertions(t, setting, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDefaultDailySummaryUserSettings(t *testing.T) {
	dailySummaryUserSettings := DefaultDailySummaryUserSettings()

	require.Equal(t, "8:00AM", dailySummaryUserSettings.PostTime)
	require.Equal(t, "Eastern Standard Time", dailySummaryUserSettings.Timezone)
	require.Equal(t, false, dailySummaryUserSettings.Enable)
}

func TestSetPanelPostID(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error storing panel postID",
			setup: func() {
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store panel postID"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to store panel postID")
			},
		},
		{
			name: "Successful Stored panel postID",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.SetPanelPostID("mockUserID", "mockPostID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestGetPanelPostID(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, string, error)
	}{
		{
			name: "Error loading panel postID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error loading panel postID"}).Times(1)
			},
			assertions: func(t *testing.T, panelPostID string, err error) {
				require.Error(t, err)
				require.Equal(t, "", panelPostID)
				require.EqualError(t, err, "failed plugin KVGet: Error loading panel postID")
			},
		},
		{
			name: "Success loading panel postID",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte("mockPostID"), nil).Times(1)
			},
			assertions: func(t *testing.T, panelPostID string, err error) {
				require.NoError(t, err)
				require.Equal(t, "mockPostID", panelPostID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			panelPostID, err := store.GetPanelPostID("mockSubscriptionID")

			tt.assertions(t, panelPostID, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeletePanelPostID(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error deleting panel post id",
			setup: func() {
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(&model.AppError{Message: "Failed to delete panel post id"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to delete panel post id")
			},
		},
		{
			name: "Successful Delete panel post id",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.DeletePanelPostID("mockUserID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}
