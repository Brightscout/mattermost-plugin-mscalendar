// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
)

func (c *client) DeleteCalendar(remoteUserID string, calID string) error {
	if !c.tokenHelpers.CheckUserStatus() {
		c.Logger.Warnf(LogUserInactive)
		return errors.New(ErrorUserInactive)
	}
	err := c.rbuilder.Users().ID(remoteUserID).Calendars().ID(calID).Request().Delete(c.ctx)
	if err != nil {
		c.tokenHelpers.ChangeUserStatus(err)
		return errors.Wrap(err, "msgraph DeleteCalendar")
	}
	c.Logger.With(bot.LogContext{}).Infof("msgraph: DeleteCalendar deleted calendar `%v`.", calID)
	return nil
}
