// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
)

// CreateEvent creates a calendar event
func (c *client) CreateEvent(remoteUserID string, in *serializer.Event) (*serializer.Event, error) {
	var out = serializer.Event{}
	if !c.tokenHelpers.CheckUserStatus(c.mattermostUserID) {
		c.Logger.Warnf(LogUserInactive, c.mattermostUserID)
		return nil, errors.New(ErrorUserInactive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Events().Request().JSONRequest(c.ctx, http.MethodPost, "", &in, &out)
	if err != nil {
		c.tokenHelpers.ChangeUserStatus(err, c.mattermostUserID)
		return nil, errors.Wrap(err, "msgraph CreateEvent")
	}
	return &out, nil
}
