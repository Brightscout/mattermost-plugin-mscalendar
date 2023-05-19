// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-mscalendar/server/store (interfaces: Store)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	store "github.com/mattermost/mattermost-plugin-mscalendar/server/store"
	bot "github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
	oauth2 "golang.org/x/oauth2"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// DeleteCurrentStep mocks base method
func (m *MockStore) DeleteCurrentStep(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCurrentStep", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCurrentStep indicates an expected call of DeleteCurrentStep
func (mr *MockStoreMockRecorder) DeleteCurrentStep(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCurrentStep", reflect.TypeOf((*MockStore)(nil).DeleteCurrentStep), arg0)
}

// DeletePanelPostID mocks base method
func (m *MockStore) DeletePanelPostID(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePanelPostID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePanelPostID indicates an expected call of DeletePanelPostID
func (mr *MockStoreMockRecorder) DeletePanelPostID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePanelPostID", reflect.TypeOf((*MockStore)(nil).DeletePanelPostID), arg0)
}

// DeleteUser mocks base method
func (m *MockStore) DeleteUser(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockStoreMockRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0)
}

// DeleteUserEvent mocks base method
func (m *MockStore) DeleteUserEvent(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserEvent indicates an expected call of DeleteUserEvent
func (mr *MockStoreMockRecorder) DeleteUserEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserEvent", reflect.TypeOf((*MockStore)(nil).DeleteUserEvent), arg0, arg1)
}

// DeleteUserFromIndex mocks base method
func (m *MockStore) DeleteUserFromIndex(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserFromIndex", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserFromIndex indicates an expected call of DeleteUserFromIndex
func (mr *MockStoreMockRecorder) DeleteUserFromIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserFromIndex", reflect.TypeOf((*MockStore)(nil).DeleteUserFromIndex), arg0)
}

// DeleteUserSubscription mocks base method
func (m *MockStore) DeleteUserSubscription(arg0 *store.User, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSubscription", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSubscription indicates an expected call of DeleteUserSubscription
func (mr *MockStoreMockRecorder) DeleteUserSubscription(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSubscription", reflect.TypeOf((*MockStore)(nil).DeleteUserSubscription), arg0, arg1)
}

// DeleteUserWelcomePost mocks base method
func (m *MockStore) DeleteUserWelcomePost(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserWelcomePost", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUserWelcomePost indicates an expected call of DeleteUserWelcomePost
func (mr *MockStoreMockRecorder) DeleteUserWelcomePost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserWelcomePost", reflect.TypeOf((*MockStore)(nil).DeleteUserWelcomePost), arg0)
}

// GetCurrentStep mocks base method
func (m *MockStore) GetCurrentStep(arg0 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentStep", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentStep indicates an expected call of GetCurrentStep
func (mr *MockStoreMockRecorder) GetCurrentStep(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentStep", reflect.TypeOf((*MockStore)(nil).GetCurrentStep), arg0)
}

// GetPanelPostID mocks base method
func (m *MockStore) GetPanelPostID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPanelPostID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPanelPostID indicates an expected call of GetPanelPostID
func (mr *MockStoreMockRecorder) GetPanelPostID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPanelPostID", reflect.TypeOf((*MockStore)(nil).GetPanelPostID), arg0)
}

// GetPostID mocks base method
func (m *MockStore) GetPostID(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostID", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostID indicates an expected call of GetPostID
func (mr *MockStoreMockRecorder) GetPostID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostID", reflect.TypeOf((*MockStore)(nil).GetPostID), arg0, arg1)
}

// GetSetting mocks base method
func (m *MockStore) GetSetting(arg0, arg1 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSetting", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSetting indicates an expected call of GetSetting
func (mr *MockStoreMockRecorder) GetSetting(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSetting", reflect.TypeOf((*MockStore)(nil).GetSetting), arg0, arg1)
}

// LoadMattermostUserID mocks base method
func (m *MockStore) LoadMattermostUserID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMattermostUserID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadMattermostUserID indicates an expected call of LoadMattermostUserID
func (mr *MockStoreMockRecorder) LoadMattermostUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMattermostUserID", reflect.TypeOf((*MockStore)(nil).LoadMattermostUserID), arg0)
}

// LoadSubscription mocks base method
func (m *MockStore) LoadSubscription(arg0 string) (*store.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadSubscription", arg0)
	ret0, _ := ret[0].(*store.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadSubscription indicates an expected call of LoadSubscription
func (mr *MockStoreMockRecorder) LoadSubscription(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadSubscription", reflect.TypeOf((*MockStore)(nil).LoadSubscription), arg0)
}

// LoadUser mocks base method
func (m *MockStore) LoadUser(arg0 string) (*store.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUser", arg0)
	ret0, _ := ret[0].(*store.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUser indicates an expected call of LoadUser
func (mr *MockStoreMockRecorder) LoadUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUser", reflect.TypeOf((*MockStore)(nil).LoadUser), arg0)
}

// LoadUserEvent mocks base method
func (m *MockStore) LoadUserEvent(arg0, arg1 string) (*store.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserEvent", arg0, arg1)
	ret0, _ := ret[0].(*store.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserEvent indicates an expected call of LoadUserEvent
func (mr *MockStoreMockRecorder) LoadUserEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserEvent", reflect.TypeOf((*MockStore)(nil).LoadUserEvent), arg0, arg1)
}

// LoadUserFromIndex mocks base method
func (m *MockStore) LoadUserFromIndex(arg0 string) (*store.UserShort, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserFromIndex", arg0)
	ret0, _ := ret[0].(*store.UserShort)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserFromIndex indicates an expected call of LoadUserFromIndex
func (mr *MockStoreMockRecorder) LoadUserFromIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserFromIndex", reflect.TypeOf((*MockStore)(nil).LoadUserFromIndex), arg0)
}

// LoadUserIndex mocks base method
func (m *MockStore) LoadUserIndex() (store.UserIndex, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserIndex")
	ret0, _ := ret[0].(store.UserIndex)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserIndex indicates an expected call of LoadUserIndex
func (mr *MockStoreMockRecorder) LoadUserIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserIndex", reflect.TypeOf((*MockStore)(nil).LoadUserIndex))
}

// LoadUserWelcomePost mocks base method
func (m *MockStore) LoadUserWelcomePost(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserWelcomePost", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserWelcomePost indicates an expected call of LoadUserWelcomePost
func (mr *MockStoreMockRecorder) LoadUserWelcomePost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserWelcomePost", reflect.TypeOf((*MockStore)(nil).LoadUserWelcomePost), arg0)
}

// MakeChangeUserStatus mocks base method
func (m *MockStore) MakeChangeUserStatus(arg0 error, arg1 bot.Logger, arg2 string, arg3 bot.Poster) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MakeChangeUserStatus", arg0, arg1, arg2, arg3)
}

// MakeChangeUserStatus indicates an expected call of MakeChangeUserStatus
func (mr *MockStoreMockRecorder) MakeChangeUserStatus(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeChangeUserStatus", reflect.TypeOf((*MockStore)(nil).MakeChangeUserStatus), arg0, arg1, arg2, arg3)
}

// MakeCheckUserStatus mocks base method
func (m *MockStore) MakeCheckUserStatus(arg0 bot.Logger, arg1 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeCheckUserStatus", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// MakeCheckUserStatus indicates an expected call of MakeCheckUserStatus
func (mr *MockStoreMockRecorder) MakeCheckUserStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeCheckUserStatus", reflect.TypeOf((*MockStore)(nil).MakeCheckUserStatus), arg0, arg1)
}

// ModifyUserIndex mocks base method
func (m *MockStore) ModifyUserIndex(arg0 func(store.UserIndex) (store.UserIndex, error)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifyUserIndex", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ModifyUserIndex indicates an expected call of ModifyUserIndex
func (mr *MockStoreMockRecorder) ModifyUserIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifyUserIndex", reflect.TypeOf((*MockStore)(nil).ModifyUserIndex), arg0)
}

// RefreshAndStoreToken mocks base method
func (m *MockStore) RefreshAndStoreToken(arg0 *oauth2.Token, arg1 *oauth2.Config, arg2 string) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshAndStoreToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshAndStoreToken indicates an expected call of RefreshAndStoreToken
func (mr *MockStoreMockRecorder) RefreshAndStoreToken(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshAndStoreToken", reflect.TypeOf((*MockStore)(nil).RefreshAndStoreToken), arg0, arg1, arg2)
}

// RemovePostID mocks base method
func (m *MockStore) RemovePostID(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePostID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePostID indicates an expected call of RemovePostID
func (mr *MockStoreMockRecorder) RemovePostID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePostID", reflect.TypeOf((*MockStore)(nil).RemovePostID), arg0, arg1)
}

// SetCurrentStep mocks base method
func (m *MockStore) SetCurrentStep(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrentStep", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrentStep indicates an expected call of SetCurrentStep
func (mr *MockStoreMockRecorder) SetCurrentStep(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrentStep", reflect.TypeOf((*MockStore)(nil).SetCurrentStep), arg0, arg1)
}

// SetPanelPostID mocks base method
func (m *MockStore) SetPanelPostID(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPanelPostID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPanelPostID indicates an expected call of SetPanelPostID
func (mr *MockStoreMockRecorder) SetPanelPostID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPanelPostID", reflect.TypeOf((*MockStore)(nil).SetPanelPostID), arg0, arg1)
}

// SetPostID mocks base method
func (m *MockStore) SetPostID(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPostID", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPostID indicates an expected call of SetPostID
func (mr *MockStoreMockRecorder) SetPostID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPostID", reflect.TypeOf((*MockStore)(nil).SetPostID), arg0, arg1, arg2)
}

// SetProperty mocks base method
func (m *MockStore) SetProperty(arg0, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProperty", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetProperty indicates an expected call of SetProperty
func (mr *MockStoreMockRecorder) SetProperty(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProperty", reflect.TypeOf((*MockStore)(nil).SetProperty), arg0, arg1, arg2)
}

// SetSetting mocks base method
func (m *MockStore) SetSetting(arg0, arg1 string, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSetting", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSetting indicates an expected call of SetSetting
func (mr *MockStoreMockRecorder) SetSetting(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSetting", reflect.TypeOf((*MockStore)(nil).SetSetting), arg0, arg1, arg2)
}

// StoreOAuth2State mocks base method
func (m *MockStore) StoreOAuth2State(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreOAuth2State", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreOAuth2State indicates an expected call of StoreOAuth2State
func (mr *MockStoreMockRecorder) StoreOAuth2State(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreOAuth2State", reflect.TypeOf((*MockStore)(nil).StoreOAuth2State), arg0)
}

// StoreUser mocks base method
func (m *MockStore) StoreUser(arg0 *store.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUser indicates an expected call of StoreUser
func (mr *MockStoreMockRecorder) StoreUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUser", reflect.TypeOf((*MockStore)(nil).StoreUser), arg0)
}

// StoreUserActiveEvents mocks base method
func (m *MockStore) StoreUserActiveEvents(arg0 string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserActiveEvents", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserActiveEvents indicates an expected call of StoreUserActiveEvents
func (mr *MockStoreMockRecorder) StoreUserActiveEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserActiveEvents", reflect.TypeOf((*MockStore)(nil).StoreUserActiveEvents), arg0, arg1)
}

// StoreUserEvent mocks base method
func (m *MockStore) StoreUserEvent(arg0 string, arg1 *store.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserEvent indicates an expected call of StoreUserEvent
func (mr *MockStoreMockRecorder) StoreUserEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserEvent", reflect.TypeOf((*MockStore)(nil).StoreUserEvent), arg0, arg1)
}

// StoreUserInIndex mocks base method
func (m *MockStore) StoreUserInIndex(arg0 *store.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserInIndex", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserInIndex indicates an expected call of StoreUserInIndex
func (mr *MockStoreMockRecorder) StoreUserInIndex(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserInIndex", reflect.TypeOf((*MockStore)(nil).StoreUserInIndex), arg0)
}

// StoreUserSubscription mocks base method
func (m *MockStore) StoreUserSubscription(arg0 *store.User, arg1 *store.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserSubscription", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserSubscription indicates an expected call of StoreUserSubscription
func (mr *MockStoreMockRecorder) StoreUserSubscription(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserSubscription", reflect.TypeOf((*MockStore)(nil).StoreUserSubscription), arg0, arg1)
}

// StoreUserWelcomePost mocks base method
func (m *MockStore) StoreUserWelcomePost(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUserWelcomePost", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUserWelcomePost indicates an expected call of StoreUserWelcomePost
func (mr *MockStoreMockRecorder) StoreUserWelcomePost(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUserWelcomePost", reflect.TypeOf((*MockStore)(nil).StoreUserWelcomePost), arg0, arg1)
}

// VerifyOAuth2State mocks base method
func (m *MockStore) VerifyOAuth2State(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyOAuth2State", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyOAuth2State indicates an expected call of VerifyOAuth2State
func (mr *MockStoreMockRecorder) VerifyOAuth2State(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyOAuth2State", reflect.TypeOf((*MockStore)(nil).VerifyOAuth2State), arg0)
}
