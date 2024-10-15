package store

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestLoadUserEvent(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, *Event, error)
	}{
		{
			name: "Error loading event",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Event not found"}).Times(1)
			},
			assertions: func(t *testing.T, event *Event, err error) {
				require.Nil(t, event)
				require.EqualError(t, err, "failed plugin KVGet: Event not found")
			},
		},
		{
			name: "Successful Load",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"PluginVersion":"1.0","Remote":{"ID":"mockRemoteID"}}`), nil).Times(1)
			},
			assertions: func(t *testing.T, event *Event, err error) {
				require.NoError(t, err)
				require.Equal(t, "1.0", event.PluginVersion)
				require.Equal(t, "mockRemoteID", event.Remote.ID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			event, err := store.LoadUserEvent("mockUserID", "mockEventID")

			tt.assertions(t, event, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestAddLinkedChannelToEvent(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading event metadata",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Metadata not found"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Metadata not found")
			},
		},
		{
			name: "Successful addition of linked channel",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return(nil, nil).Times(1)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			err := store.AddLinkedChannelToEvent("mockEventID", "mockChannelID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteLinkedChannelFromEvent(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading event metadata",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Metadata not found"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Metadata not found")
			},
		},
		{
			name: "Channel ID not present",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"LinkedChannelIDs":{"otherChannelID":{}}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Error storing updated metadata",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"LinkedChannelIDs":{"mockChannelID":{}}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(&model.AppError{Message: "Failed to store metadata"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to store metadata")
			},
		},
		{
			name: "Successful deletion of linked channel",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"LinkedChannelIDs":{"mockChannelID":{}}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			err := store.DeleteLinkedChannelFromEvent("mockEventID", "mockChannelID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreEventMetadata(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error storing event metadata",
			setup: func() {
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(&model.AppError{Message: "Failed to store metadata"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "error storing event metadata")
			},
		},
		{
			name: "Successful store of event metadata",
			setup: func() {
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			eventMeta := &EventMetadata{
				LinkedChannelIDs: map[string]struct{}{
					"mockChannelID": {},
				},
			}
			err := store.StoreEventMetadata("mockEventID", eventMeta)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestLoadEventMetadata(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, *EventMetadata, error)
	}{
		{
			name: "Error loading event metadata",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Failed to load event metadata"}).Times(1)
			},
			assertions: func(t *testing.T, eventMeta *EventMetadata, err error) {
				require.Nil(t, eventMeta)
				require.ErrorContains(t, err, "Failed to load event metadata")
			},
		},
		{
			name: "Successful load of event metadata",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte(`{"LinkedChannelIDs":{"mockChannelID":{}}}`), nil).Times(1)
			},
			assertions: func(t *testing.T, eventMeta *EventMetadata, err error) {
				require.NoError(t, err)
				require.Contains(t, eventMeta.LinkedChannelIDs, "mockChannelID")
			},
		},
		{
			name: "Event metadata not found",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, nil).Times(1)
			},
			assertions: func(t *testing.T, eventMeta *EventMetadata, err error) {
				require.ErrorContains(t, err, "not found")
				require.Nil(t, eventMeta)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			eventMeta, err := store.LoadEventMetadata("mockEventID")

			tt.assertions(t, eventMeta, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteEventMetadata(t *testing.T) {
	mockAPI, store, _, _, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error deleting event metadata",
			setup: func() {
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(&model.AppError{Message: "Failed to delete event metadata"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Failed to delete event metadata")
			},
		},
		{
			name: "Successful deletion of event metadata",
			setup: func() {
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			err := store.DeleteEventMetadata("mockEventID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserEvent(t *testing.T) {
	mockAPI, store, mockLogger, mockLoggerWith, _ := MockStoreSetup(t)
	mockEvent := &Event{Remote: &remote.Event{ICalUID: "mockICalUID", ID: "mockEventID"}}
	mockUserID := "user1"

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Store expired event",
			setup: func() {
				mockEvent.Remote.End = &remote.DateTime{DateTime: "2006-01-02T15:04:05"}
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Error storing user event",
			setup: func() {
				mockEvent.Remote.End = remote.NewDateTime(time.Now(), "UTC")
				mockAPI.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("int64")).Return(&model.AppError{Message: "Failed to store user event"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to store user event")
			},
		},
		{
			name: "Successful store user event",
			setup: func() {
				mockEvent.Remote.End = remote.NewDateTime(time.Now(), "UTC")
				mockAPI.On("KVSetWithExpiry", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("int64")).Return(nil).Times(1)
				mockLogger.EXPECT().With(gomock.Any()).Return(mockLoggerWith).Times(1)
				mockLoggerWith.EXPECT().Debugf("store: stored user event.").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			err := store.StoreUserEvent(mockUserID, mockEvent)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteUserEvent(t *testing.T) {
	mockAPI, store, mockLogger, mockLoggerWith, _ := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error deleting user event",
			setup: func() {
				mockAPI.On("KVDelete", mock.Anything).Return(&model.AppError{Message: "Failed to delete event"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Failed to delete event")
			},
		},
		{
			name: "Successful delete",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVDelete", mock.Anything).Return(nil).Times(1)
				mockLogger.EXPECT().With(gomock.Any()).Return(mockLoggerWith).Times(1)
				mockLoggerWith.EXPECT().Debugf("store: deleted event.").Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI.ExpectedCalls = nil
			tt.setup()

			err := store.DeleteUserEvent("mockUserID", "mockEventID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}
