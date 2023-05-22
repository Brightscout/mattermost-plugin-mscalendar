// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"context"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
)

type Client interface {
	MakeClient() (remote.Client, error)
	MakeSuperuserClient() (remote.Client, error)
}

func (m *mscalendar) MakeClient() (remote.Client, error) {
	err := m.Filter(withActingUserExpanded)
	if err != nil {
		return nil, err
	}

	tokenHelpers := &remote.UserTokenHelpers{
		CheckUserStatus:      m.Store.CheckUserConnected,
		ChangeUserStatus:     m.Store.DisconnectUserFromStoreIfNecessary,
		RefreshAndStoreToken: m.Store.RefreshAndStoreToken,
	}

	return m.Remote.MakeUserClient(context.Background(), m.actingUser.OAuth2Token, m.actingUser.MattermostUserID, m.Poster, tokenHelpers), nil
}

func (m *mscalendar) MakeSuperuserClient() (remote.Client, error) {
	return m.Remote.MakeSuperuserClient(context.Background())
}
