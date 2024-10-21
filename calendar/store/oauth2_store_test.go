package store

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestVerifyOAuth2State(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading state",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return(nil, &model.AppError{Message: "Error getting state"}).Times(1)
				mockAPI.On("KVDelete", mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed plugin KVGet: Error getting state")
			},
		},
		{
			name: "Invalid Oauth state",
			setup: func() {
				mockAPI.On("KVGet", mock.Anything).Return([]byte("invalidState"), nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid oauth state, please try again")
			},
		},
		{
			name: "Successfull Oauth state verification",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.Anything).Return([]byte("mockState"), nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.VerifyOAuth2State("mockState")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreOAuth2State(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func()
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading state",
			setup: func() {
				mockAPI.On("KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything).Return(&model.AppError{Message: "Error loading state"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "Error loading state")
			},
		},
		{
			name: "Successfull Oauth state verification",
			setup: func() {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := store.StoreOAuth2State("mockState")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}
