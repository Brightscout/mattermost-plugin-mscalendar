// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
)

// CreateCalendar creates a calendar
func (c *client) CreateCalendar(remoteUserID string, calIn *remote.Calendar) (*remote.Calendar, error) {
	var calOut = remote.Calendar{}
	if !c.checkUserStatus() {
		c.Logger.Warnf(LogUserInactive)
		return nil, errors.New(ErrorUserInactive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Calendars().Request().JSONRequest(c.ctx, http.MethodPost, "", &calIn, &calOut)
	if err != nil {
		c.changeUserStatus(err)
		return nil, errors.Wrap(err, "msgraph CreateCalendar")
	}
	c.Logger.With(bot.LogContext{
		"v": calOut,
	}).Infof("msgraph: CreateCalendar created the following calendar.")
	return &calOut, nil
}
