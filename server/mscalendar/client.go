// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package mscalendar

import (
	"context"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/store"
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

	return m.Remote.MakeClient(context.Background(), m.actingUser.OAuth2Token, m.Store, m.actingUser.MattermostUserID, store.GetCheckUserStatus(m.Store, m.Logger, m.actingUser.MattermostUserID), store.GetChangeUserStatus(m.Store, m.Logger, m.actingUser.MattermostUserID, m.Poster)), nil
}

func (m *mscalendar) MakeSuperuserClient() (remote.Client, error) {
	return m.Remote.MakeSuperuserClient(context.Background())
}
