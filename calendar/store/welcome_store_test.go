package store

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type MockPluginAPI struct {
	plugin.API
	mock.Mock
}

func (m *MockPluginAPI) KVGet(key string) ([]byte, *model.AppError) {
	args := m.Called(key)
	data, _ := args.Get(0).([]byte)
	if err := args.Get(1); err != nil {
		return nil, err.(*model.AppError)
	}
	return data, nil
}

func (m *MockPluginAPI) KVSet(key string, data []byte) *model.AppError {
	args := m.Called(key, data)
	if err := args.Get(0); err != nil {
		return err.(*model.AppError)
	}
	return nil
}

func (m *MockPluginAPI) KVDelete(key string) *model.AppError {
	args := m.Called(key)
	if err := args.Get(0); err != nil {
		return err.(*model.AppError)
	}
	return nil
}

func (m *MockPluginAPI) KVSetWithOptions(key string, value []byte, options model.PluginKVSetOptions) (bool, *model.AppError) {
	args := m.Called(key, value, options)

	success := args.Bool(0)
	if err := args.Get(1); err != nil {
		return success, err.(*model.AppError)
	}
	return success, nil
}

func TestLoadUserWelcomePost(t *testing.T) {
	mockAPI, store := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func(*MockPluginAPI)
		assertions func(*testing.T, string, error)
	}{
		{
			name: "Error loading user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"})
			},
			assertions: func(t *testing.T, resp string, err error) {
				require.Equal(t, resp, "")
				require.EqualError(t, err, "failed plugin KVGet: KVGet failed")
			},
		},
		{
			name: "Success loading user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`"mockPostID"`), nil)
			},
			assertions: func(t *testing.T, resp string, err error) {
				require.NoError(t, err)
				require.Equal(t, "mockPostID", resp)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			resp, err := store.LoadUserWelcomePost("mockMMUserID")

			tt.assertions(t, resp, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserWelcomePost(t *testing.T) {
	mockAPI, store := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func(*MockPluginAPI)
		assertions func(*testing.T, error)
	}{
		{
			name: "Error storing user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(&model.AppError{Message: "KVSet failed"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContainsf(t, err, "failed plugin KVSet (ttl: 0s)", `"mockMMUserID": KVSet failed`)
			},
		},
		{
			name: "Success storing user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.StoreUserWelcomePost("mockMMUserID", "mockPostID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteUserWelcomePost(t *testing.T) {
	mockAPI, store := MockStoreSetup(t)

	tests := []struct {
		name       string
		setup      func(*MockPluginAPI)
		assertions func(*testing.T, string, error)
	}{
		{
			name: "Error deleting user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`"mockPostID"`), nil)
				mockAPI.On("KVDelete", mock.Anything, mock.Anything).Return(&model.AppError{Message: "KVDelete failed"})
			},
			assertions: func(t *testing.T, resp string, err error) {
				require.Equal(t, "", resp)
				require.ErrorContains(t, err, "failed plugin KVdelete", "KVDelete failed")
			},
		},
		{
			name: "Success deleting user welcome post",
			setup: func(mockAPI *MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`"mockPostID"`), nil)
				mockAPI.On("KVDelete", mock.Anything, mock.Anything).Return(nil)
			},
			assertions: func(t *testing.T, resp string, err error) {
				require.Equal(t, "mockPostID", resp)
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			resp, err := store.DeleteUserWelcomePost("mockMMUserID")

			tt.assertions(t, resp, err)

			mockAPI.AssertExpectations(t)
		})
	}
}
