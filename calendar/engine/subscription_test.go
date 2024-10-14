package engine

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/calendar/store"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestCreateMyEventSubscription(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()
	expectedSub := GetMockSubscription()

	tests := []struct {
		name      string
		setupMock func()
		assertion func(sub *store.Subscription, err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(sub *store.Subscription, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error creating subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mockClient.EXPECT().CreateMySubscription(gomock.Any(), "testActingUserRemoteID").Return(nil, errors.New("error creating subscription"))
			},
			assertion: func(sub *store.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error creating subscription")
			},
		},
		{
			name: "subscription created successfully",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mockClient.EXPECT().CreateMySubscription(gomock.Any(), "testActingUserRemoteID").Return(&remote.Subscription{}, nil)
				mockStore.EXPECT().StoreUserSubscription(mscalendar.actingUser.User, expectedSub)
			},
			assertion: func(sub *store.Subscription, err error) {
				require.NoError(t, err)
				require.Equal(t, expectedSub, sub)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			sub, err := mscalendar.CreateMyEventSubscription()

			tt.assertion(sub, err)
		})
	}
}

func TestLoadMyEventSubscription(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()
	expectedSubscription := GetMockSubscription()

	tests := []struct {
		name      string
		setupMock func()
		assertion func(sub *store.Subscription, err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(sub *store.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name: "error loading subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(nil, errors.New("error loading subscription")).Times(1)
			},
			assertion: func(sub *store.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error loading subscription")
			},
		},
		{
			name: "subscription loaded successfully",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(expectedSubscription, nil).Times(1)
			},
			assertion: func(sub *store.Subscription, err error) {
				require.NoError(t, err)
				require.Equal(t, expectedSubscription, sub)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			sub, err := mscalendar.LoadMyEventSubscription()

			tt.assertion(sub, err)
		})
	}
}

func TestListRemoteSubscriptions(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name      string
		setupMock func()
		assertion func(subs []*remote.Subscription, err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(subs []*remote.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error withClient in ListRemoteSubscriptions: It looks like your Mattermost account is not connected to testDisplayName. Please connect your account using `/testCommandTrigger connect`.: error filtering user")
			},
		},
		{
			name: "error listing subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockClient.EXPECT().ListSubscriptions().Return(nil, errors.New("error listing subscriptions"))
			},
			assertion: func(subs []*remote.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error listing subscriptions")
			},
		},
		{
			name: "subscriptions listed successfully",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockClient.EXPECT().ListSubscriptions().Return([]*remote.Subscription{{ID: "sub1"}, {ID: "sub2"}}, nil)
			},
			assertion: func(subs []*remote.Subscription, err error) {
				require.NoError(t, err)
				require.Equal(t, []*remote.Subscription{{ID: "sub1"}, {ID: "sub2"}}, subs)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			subs, err := mscalendar.ListRemoteSubscriptions()

			tt.assertion(subs, err)
		})
	}
}

func TestRenewMyEventSubscription(t *testing.T) {
	mscalendar, mockStore, _, mockRemote, mockPluginAPI, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name      string
		setupMock func()
		assertion func(subs *store.Subscription, err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(subs *store.Subscription, err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "no subscriptions present",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil)
				mockRemote.EXPECT().MakeClient(gomock.Any(), nil)
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = ""
			},
			assertion: func(subs *store.Subscription, err error) {
				require.NoError(t, err)
				require.Nil(t, subs)
			},
		},
		{
			name: "error loading subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil)
				mockRemote.EXPECT().MakeClient(gomock.Any(), nil)
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(nil, errors.New("some error occurred while loading subscription"))
			},
			assertion: func(subs *store.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error loading subscription: some error occurred while loading subscription")
			},
		},
		{
			name: "error renewing subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil)
				mockClient.EXPECT().RenewSubscription(gomock.Any(), "testActingUserRemoteID", &remote.Subscription{}).Return(nil, errors.New("The object was not found")).Times(1)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "testEventSubscriptionID").Return(errors.New("error deleting subscription")).Times(1)
			},
			assertion: func(subs *store.Subscription, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error deleting subscription")
			},
		},
		{
			name: "successfully renew event subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil).Times(2)
				mockClient.EXPECT().RenewSubscription(gomock.Any(), "testActingUserRemoteID", &remote.Subscription{}).Return(&remote.Subscription{}, nil).Times(1)
				mockStore.EXPECT().StoreUserSubscription(gomock.Any(), gomock.Any()).Return(nil)
			},
			assertion: func(subs *store.Subscription, err error) {
				require.NoError(t, err)
				require.Equal(t, &store.Subscription{Remote: &remote.Subscription{}}, subs)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			subs, err := mscalendar.RenewMyEventSubscription()

			tt.assertion(subs, err)
		})
	}
}

func TestDeleteMyEventSubscription(t *testing.T) {
	mscalendar, mockStore, _, _, mockPluginAPI, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name      string
		setupMock func()
		assertion func(err error)
	}{
		{
			name: "error filtering with user",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name: "error loading subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(nil, errors.New("some error occurred while loading subscription")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "error loading subscription: some error occurred while loading subscription")
			},
		},
		{
			name: "error deleting subscription in DeleteOrphanedSubscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil).Times(1)
				mscalendar.client = mockClient
				mockClient.EXPECT().DeleteSubscription(&remote.Subscription{}).Return(errors.New("some error occured")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed to delete subscription : some error occured")
			},
		},
		{
			name: "error deleting user subscription",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil).Times(1)
				mscalendar.client = mockClient
				mockClient.EXPECT().DeleteSubscription(&remote.Subscription{}).Return(nil).Times(1)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "testEventSubscriptionID").Return(errors.New("error deleting user subscription"))
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed to delete subscription testEventSubscriptionID: error deleting user subscription")
			},
		},
		{
			name: "event subscription deleted successfully",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mscalendar.actingUser.Settings.EventSubscriptionID = "testEventSubscriptionID"
				mockPluginAPI.EXPECT().GetMattermostUser("testActingUserID").Return(&model.User{}, nil).Times(1)
				mockStore.EXPECT().LoadSubscription("testEventSubscriptionID").Return(&store.Subscription{Remote: &remote.Subscription{}}, nil).Times(1)
				mscalendar.client = mockClient
				mockClient.EXPECT().DeleteSubscription(&remote.Subscription{}).Return(nil).Times(1)
				mockStore.EXPECT().DeleteUserSubscription(gomock.Any(), "testEventSubscriptionID").Return(nil)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := mscalendar.DeleteMyEventSubscription()

			tt.assertion(err)
		})
	}
}

func TestDeleteOrphanedSubscription(t *testing.T) {
	mscalendar, mockStore, _, _, _, mockClient, _ := MockSetup(t)
	mockUser := GetMockUser()

	tests := []struct {
		name      string
		eventID   string
		setupMock func()
		assertion func(err error)
	}{
		{
			name:    "error filtering with user",
			eventID: "testEventID",
			setupMock: func() {
				mockUser.User = nil
				mscalendar.client = nil
				mscalendar.actingUser = &User{MattermostUserID: "testActingUserID"}
				mockStore.EXPECT().LoadUser("testActingUserID").Return(nil, errors.New("error filtering user")).Times(1)
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.ErrorContains(t, err, "error filtering user")
			},
		},
		{
			name:    "error deleting subscription",
			eventID: "testEventID",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mockClient.EXPECT().DeleteSubscription(gomock.Any()).Return(errors.New("error deleting subscription"))
			},
			assertion: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "failed to delete subscription : error deleting subscription")
			},
		},
		{
			name:    "subscription deleted sucessfully",
			eventID: "testEventID",
			setupMock: func() {
				mockUser.User = &store.User{Settings: store.Settings{}, Remote: &remote.User{ID: "testRemoteID"}}
				mscalendar.client = mockClient
				mscalendar.actingUser = &User{User: &store.User{Remote: &remote.User{ID: "testActingUserRemoteID"}}, MattermostUserID: "testActingUserID"}
				mockClient.EXPECT().DeleteSubscription(gomock.Any()).Return(nil)
			},
			assertion: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			subscription := GetMockSubscription()

			err := mscalendar.DeleteOrphanedSubscription(subscription)

			tt.assertion(err)
		})
	}
}
