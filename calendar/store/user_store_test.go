package store

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/testutil"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestLoadUser(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, *User, error)
	}{
		{
			name: "Error loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"})
			},
			assertions: func(t *testing.T, user *User, err error) {
				require.Nil(t, user)
				require.EqualError(t, err, "failed plugin KVGet: KVGet failed")
			},
		},
		{
			name: "Success loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"isCustomStatusSet": false}`), nil)
			},
			assertions: func(t *testing.T, user *User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.False(t, user.IsCustomStatusSet)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			user, err := store.LoadUser("mockMMUserID")

			tt.assertions(t, user, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestLoadMattermostUserID(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, string, error)
	}{
		{
			name: "Error loading Mattermost User ID",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "Load failed"})
			},
			assertions: func(t *testing.T, userID string, err error) {
				require.Empty(t, userID)
				require.EqualError(t, err, "failed plugin KVGet: Load failed")
			},
		},
		{
			name: "Success loading Mattermost User ID",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockUserID := []byte("mockMattermostUserID")
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(mockUserID, nil)
			},
			assertions: func(t *testing.T, userID string, err error) {
				require.NoError(t, err)
				require.Equal(t, "mockMattermostUserID", userID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			userID, err := store.LoadMattermostUserID("mockRemoteUserID")

			tt.assertions(t, userID, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestLoadUserIndex(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, UserIndex, error)
	}{
		{
			name: "Error loading UserIndex",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "Load failed"})
			},
			assertions: func(t *testing.T, userIndex UserIndex, err error) {
				require.Nil(t, userIndex)
				require.EqualError(t, err, "failed plugin KVGet: Load failed")
			},
		},
		{
			name: "Success loading UserIndex",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockUserIndexJSON := `[{"mm_username": "mockUser"}]`
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(mockUserIndexJSON), nil)
			},
			assertions: func(t *testing.T, userIndex UserIndex, err error) {
				require.NoError(t, err)
				require.Len(t, userIndex, 1)
				require.Equal(t, "mockUser", userIndex[0].MattermostUsername)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			userIndex, err := store.LoadUserIndex()

			tt.assertions(t, userIndex, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestLoadUserFromIndex(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, *UserShort, error)
	}{
		{
			name: "Error loading UserIndex",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "Load failed"})
			},
			assertions: func(t *testing.T, user *UserShort, err error) {
				require.Nil(t, user)
				require.EqualError(t, err, "failed plugin KVGet: Load failed")
			},
		},
		{
			name: "User not found in index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockUserIndexJSON := `[{"mm_id": "mockMMUserID2"}]`
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(mockUserIndexJSON), nil)
			},
			assertions: func(t *testing.T, user *UserShort, err error) {
				require.Nil(t, user)
				require.Equal(t, ErrNotFound, err)
			},
		},
		{
			name: "Success loading User from index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockUserIndexJSON := `[{"mm_id": "mockMMUserID"}]`
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(mockUserIndexJSON), nil)
			},
			assertions: func(t *testing.T, user *UserShort, err error) {
				require.NoError(t, err)
				require.Equal(t, "mockMMUserID", user.MattermostUserID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			user, err := store.LoadUserFromIndex("mockMMUserID")

			tt.assertions(t, user, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUser(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	user := &User{
		MattermostUserID:      "mockMMUserID",
		Remote:                &remote.User{ID: "mockRemoteID"},
		MattermostUsername:    "mockUser",
		MattermostDisplayName: "Mock User",
	}

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, error)
	}{
		{
			name: "Error storing user JSON",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store user"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "failed plugin KVSet", "Failed to store user")
			},
		},
		{
			name: "Error storing Mattermost User ID",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), []byte("mockMMUserID")).Return(&model.AppError{Message: "Failed to store Mattermost User ID"}).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "failed plugin KVSet", "Failed to store user")
			},
		},
		{
			name: "Success storing user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), []byte("mockMMUserID")).Return(nil)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.StoreUser(user)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"})
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "failed plugin KVGet: KVGet failed")
			},
		},
		{
			name: "Error deleting user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{}`), nil)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(&model.AppError{Message: "error deleting user"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "error deleting user")
			},
		},
		{
			name: "Error deleting mattermost user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"remote": {"id": "mockRemoteID"}}`), nil)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(&model.AppError{Message: "error deleting mattermost user"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "error deleting mattermost user")
			},
		},
		{
			name: "Error getting user details",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"remote": {"id": "mockRemoteID"}}`), nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "error getting user details"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "error getting user details")
			},
		},
		{
			name: "error storing user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"remote": {"id": "mockRemoteID"}}`), nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(1)
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
				mockAPI.On("KVSet", mock.Anything, mock.Anything).Return(&model.AppError{Message: "error storing user"})
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "error storing user")
			},
		},
		{
			name: "Success deleting user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"remote": {"id": "mockRemoteID"}}`), nil).Times(1)
				mockAPI.On("KVDelete", mock.AnythingOfType("string")).Return(nil).Times(2)
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
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

			err := store.DeleteUser("mockMMUserID")

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserInIndex(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "modification error: failed plugin KVGet: KVGet failed")
			},
		},
		{
			name: "Error unmarshalling existing user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`invalid json`), nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "modification error: invalid character 'i' looking for beginning of value")
			},
		},
		{
			name: "Error storing updated user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(false, &model.AppError{Message: "KVSet failed"}).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "problem writing value", "KVSet failed")
			},
		},
		{
			name: "Successfully update an existing user in index with matching IDs",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[{"MattermostUserID":"mockMMUserID","RemoteID":"mockRemoteID"}]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(true, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Successfully store a new user in index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(true, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Successfully update an existing user in index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[{"MattermostUserID":"mockMMUserID","RemoteID":"mockRemoteID"}]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(true, nil).Times(1)
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)
			user := &User{
				MattermostUserID:      "mockMMUserID",
				MattermostUsername:    "mockMMUsername",
				MattermostDisplayName: "mockDisplayName",
				Remote: &remote.User{
					ID:   "mockRemoteID",
					Mail: "mock@remote.com",
				},
			}

			err := store.StoreUserInIndex(user)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestDeleteUserFromIndex(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		userID     string
		assertions func(*testing.T, error)
	}{
		{
			name: "Error loading user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"}).Times(1)
			},
			userID: "mockMMUserID",
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "modification error: failed plugin KVGet: KVGet failed")
			},
		},
		{
			name: "Error unmarshalling existing user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`invalid json`), nil).Times(1)
			},
			userID: "mockMMUserID",
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "modification error: invalid character 'i' looking for beginning of value")
			},
		},
		{
			name: "User not found in index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
			},
			userID: "mockMMUserID",
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Successfully delete a user from index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).
					Return([]byte(`[{"MattermostUserID":"mockMMUserID","RemoteID":"mockRemoteID"},{"MattermostUserID":"otherUserID","RemoteID":"otherRemoteID"}]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(true, nil).Times(1)
			},
			userID: "mockMMUserID",
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Error storing updated user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[{"MattermostUserID":"mockMMUserID","RemoteID":"mockRemoteID"}]`), nil).Times(1)
				mockAPI.On("KVSetWithOptions", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(false, &model.AppError{Message: "KVSet failed"}).Times(1)
			},
			userID: "mockMMUserID",
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "problem writing value", "KVSet failed")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.DeleteUserFromIndex(tt.userID)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestSearchInUserIndex(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name       string
		setup      func(*testutil.MockPluginAPI)
		term       string
		limit      int
		assertions func(*testing.T, UserIndex, error)
	}{
		{
			name: "Error loading user index",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "KVGet failed"}).Times(1)
			},
			term:  "searchTerm",
			limit: 5,
			assertions: func(t *testing.T, result UserIndex, err error) {
				require.EqualError(t, err, "error searching user in index: failed plugin KVGet: KVGet failed")
				require.Nil(t, result)
			},
		},
		{
			name: "No matches found",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[]`), nil).Times(1)
			},
			term:  "searchTerm",
			limit: 5,
			assertions: func(t *testing.T, result UserIndex, err error) {
				require.NoError(t, err)
				require.Empty(t, result)
			},
		},
		{
			name: "Matches found within limit",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[
					{"mm_username":"user1","remote_id":"remote1","mm_id":"user1","email":"user1@example.com"},
					{"mm_username":"user2","remote_id":"remote2","mm_id":"user2","email":"user2@example.com"}
				]`), nil).Times(1)
			},
			term:  "user",
			limit: 1,
			assertions: func(t *testing.T, result UserIndex, err error) {
				require.NoError(t, err)
				require.Len(t, result, 1)
				require.Equal(t, "user1", result[0].MattermostUserID)
			},
		},
		{
			name: "Matches not found within limit despite existing matches",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[
					{"mm_username":"user1","remote_id":"remote1","mm_id":"user1","email":"user1@example.com"},
					{"mm_username":"user2","remote_id":"remote2","mm_id":"user2","email":"user2@example.com"},
					{"mm_username":"user3","remote_id":"remote3","mm_id":"user3","email":"user3@example.com"}
				]`), nil).Times(1)
			},
			term:  "nonexistent",
			limit: 2,
			assertions: func(t *testing.T, result UserIndex, err error) {
				require.NoError(t, err)
				require.Len(t, result, 0)
			},
		},
		{
			name: "Limit exceeded but only returns available matches",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`[
					{"mm_username":"user1","remote_id":"remote1","mm_id":"user1","email":"user1@example.com"},
					{"mm_username":"user2","remote_id":"remote2","mm_id":"user2","email":"user2@example.com"},
					{"mm_username":"user3","remote_id":"remote3","mm_id":"user3","email":"user3@example.com"}
				]`), nil).Times(1)
			},
			term:  "user",
			limit: 2,
			assertions: func(t *testing.T, result UserIndex, err error) {
				require.NoError(t, err)
				require.Len(t, result, 2)
				require.Equal(t, "user1", result[0].MattermostUserID)
				require.Equal(t, "user2", result[1].MattermostUserID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			result, err := store.SearchInUserIndex(tt.term, tt.limit)

			tt.assertions(t, result, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserActiveEvents(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name             string
		setup            func(*testutil.MockPluginAPI)
		mattermostUserID string
		events           []string
		assertions       func(*testing.T, error)
	}{
		{
			name: "Error loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "User not found"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			events:           []string{"event1", "event2"},
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "failed plugin KVGet: User not found")
			},
		},
		{
			name: "Error storing active events",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","active_events": []}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store events"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			events:           []string{"event1", "event2"},
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Failed to store events")
			},
		},
		{
			name: "Store active events successfully",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","active_events": []}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
			},
			mattermostUserID: "mockUserID",
			events:           []string{"event1", "event2"},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.StoreUserActiveEvents(tt.mattermostUserID, tt.events)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserLinkedEvent(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name             string
		setup            func(*testutil.MockPluginAPI)
		mattermostUserID string
		eventID          string
		channelID        string
		assertions       func(*testing.T, error)
	}{
		{
			name: "Error loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "User not found"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			eventID:          "mockEventID",
			channelID:        "mockChannelID",
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "failed plugin KVGet: User not found")
			},
		},
		{
			name: "Error storing linked event",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","channel_events": {}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store linked event"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			eventID:          "mockEventID",
			channelID:        "mockChannelID",
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Failed to store linked event")
			},
		},
		{
			name: "Store linked event successfully with empty ChannelEvents",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","channel_events": {}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
			},
			mattermostUserID: "mockUserID",
			eventID:          "mockEventID",
			channelID:        "mockChannelID",
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Store linked event successfully with existing ChannelEvents",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","channel_events": {"mockEventID": "mockChannelID"}}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
			},
			mattermostUserID: "mockUserID",
			eventID:          "event2",
			channelID:        "channel2",
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.StoreUserLinkedEvent(tt.mattermostUserID, tt.eventID, tt.channelID)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestStoreUserCustomStatusUpdates(t *testing.T) {
	mockAPI, store, _, _, _ := GetMockSetup(t)

	tests := []struct {
		name             string
		setup            func(*testutil.MockPluginAPI)
		mattermostUserID string
		value            bool
		assertions       func(*testing.T, error)
	}{
		{
			name: "Error loading user",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return(nil, &model.AppError{Message: "User not found"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			value:            true,
			assertions: func(t *testing.T, err error) {
				require.EqualError(t, err, "failed plugin KVGet: User not found")
			},
		},
		{
			name: "Error storing custom status update",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","is_custom_status_set": false}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(&model.AppError{Message: "Failed to store custom status"}).Times(1)
			},
			mattermostUserID: "mockUserID",
			value:            true,
			assertions: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "Failed to store custom status")
			},
		},
		{
			name: "Store custom status update successfully",
			setup: func(mockAPI *testutil.MockPluginAPI) {
				mockAPI.ExpectedCalls = nil
				mockAPI.On("KVGet", mock.AnythingOfType("string")).Return([]byte(`{"mm_id":"mockUserID","is_custom_status_set": false}`), nil).Times(1)
				mockAPI.On("KVSet", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(1)
			},
			mattermostUserID: "mockUserID",
			value:            true,
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockAPI)

			err := store.StoreUserCustomStatusUpdates(tt.mattermostUserID, tt.value)

			tt.assertions(t, err)

			mockAPI.AssertExpectations(t)
		})
	}
}

func TestUserIndex_ByMattermostID(t *testing.T) {
	index := UserIndex{
		&UserShort{MattermostUserID: "user1"},
		&UserShort{MattermostUserID: "user2"},
	}

	result := index.ByMattermostID()
	require.Len(t, result, 2)
	require.Contains(t, result, "user1")
	require.Contains(t, result, "user2")
}

func TestUserIndex_ByRemoteID(t *testing.T) {
	index := UserIndex{
		&UserShort{RemoteID: "remote1"},
		&UserShort{RemoteID: "remote2"},
	}

	result := index.ByRemoteID()
	require.Len(t, result, 2)
	require.Contains(t, result, "remote1")
	require.Contains(t, result, "remote2")
}

func TestUserIndex_ByEmail(t *testing.T) {
	index := UserIndex{
		&UserShort{Email: "user1@example.com"},
		&UserShort{Email: "user2@example.com"},
	}

	result := index.ByEmail()
	require.Len(t, result, 2)
	require.Contains(t, result, "user1@example.com")
	require.Contains(t, result, "user2@example.com")
}

func TestUserIndex_GetMattermostUserIDs(t *testing.T) {
	index := UserIndex{
		&UserShort{MattermostUserID: "user1"},
		&UserShort{MattermostUserID: "user2"},
	}

	result := index.GetMattermostUserIDs()
	require.Len(t, result, 2)
	require.Contains(t, result, "user1")
	require.Contains(t, result, "user2")
}

func TestIsConfiguredForStatusUpdates(t *testing.T) {
	tests := []struct {
		name           string
		settings       Settings
		expectedResult bool
	}{
		{
			name: "UpdateStatusFromOptions is AwayStatusOption",
			settings: Settings{
				UpdateStatusFromOptions: AwayStatusOption,
			},
			expectedResult: true,
		},
		{
			name: "UpdateStatusFromOptions is DNDStatusOption",
			settings: Settings{
				UpdateStatusFromOptions: DNDStatusOption,
			},
			expectedResult: true,
		},
		{
			name: "UpdateStatusFromOptions is empty, UpdateStatus is true, notifications during meeting are false",
			settings: Settings{
				UpdateStatusFromOptions:           "",
				UpdateStatus:                      true,
				ReceiveNotificationsDuringMeeting: false,
			},
			expectedResult: true,
		},
		{
			name: "UpdateStatusFromOptions is empty, UpdateStatus is true, notifications during meeting are true",
			settings: Settings{
				UpdateStatusFromOptions:           "",
				UpdateStatus:                      true,
				ReceiveNotificationsDuringMeeting: true,
			},
			expectedResult: true,
		},
		{
			name: "UpdateStatusFromOptions is empty, UpdateStatus is false",
			settings: Settings{
				UpdateStatusFromOptions: "",
				UpdateStatus:            false,
			},
			expectedResult: false,
		},
		{
			name: "UpdateStatusFromOptions is not configured",
			settings: Settings{
				UpdateStatusFromOptions: "other",
			},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Settings: tt.settings,
			}

			result := user.IsConfiguredForStatusUpdates()
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestIsConfiguredForCustomStatusUpdates(t *testing.T) {
	tests := []struct {
		name           string
		settings       Settings
		expectedResult bool
	}{
		{
			name: "Custom status is configured",
			settings: Settings{
				SetCustomStatus: true,
			},
			expectedResult: true,
		},
		{
			name: "Custom status is not configured",
			settings: Settings{
				SetCustomStatus: false,
			},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Settings: tt.settings,
			}

			result := user.IsConfiguredForCustomStatusUpdates()
			require.Equal(t, tt.expectedResult, result, "Expected %v but got %v", tt.expectedResult, result)
		})
	}
}
