// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"net/http"

	"github.com/pkg/errors"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
)

func (c *client) GetEvent(remoteUserID, eventID string) (*serializer.Event, error) {
	e := &serializer.Event{}
	if !c.CheckUserStatus() {
		c.Logger.Warnf(LogUserInActive, c.mattermostUserID)
		return nil, errors.New(ErrorUserInActive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Request().JSONRequest(
		c.ctx, http.MethodGet, "", nil, &e)
	if err != nil {
		c.ChangeUserStatus(err)
		return nil, errors.Wrap(err, "msgraph GetEvent")
	}
	return e, nil
}

func (c *client) AcceptEvent(remoteUserID, eventID string) error {
	dummy := &msgraph.EventAcceptRequestParameter{}
	if !c.CheckUserStatus() {
		c.Logger.Warnf(LogUserInActive, c.mattermostUserID)
		return errors.New(ErrorUserInActive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Accept(dummy).Request().Post(c.ctx)
	if err != nil {
		c.ChangeUserStatus(err)
		return errors.Wrap(err, "msgraph Accept Event")
	}

	return nil
}

func (c *client) DeclineEvent(remoteUserID, eventID string) error {
	dummy := &msgraph.EventDeclineRequestParameter{}
	if !c.CheckUserStatus() {
		c.Logger.Warnf(LogUserInActive, c.mattermostUserID)
		return errors.New(ErrorUserInActive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).Decline(dummy).Request().Post(c.ctx)
	if err != nil {
		c.ChangeUserStatus(err)
		return errors.Wrap(err, "msgraph DeclineEvent")
	}
	return nil
}

func (c *client) TentativelyAcceptEvent(remoteUserID, eventID string) error {
	dummy := &msgraph.EventTentativelyAcceptRequestParameter{}
	if !c.CheckUserStatus() {
		c.Logger.Warnf(LogUserInActive, c.mattermostUserID)
		return errors.New(ErrorUserInActive)
	}

	err := c.rbuilder.Users().ID(remoteUserID).Events().ID(eventID).TentativelyAccept(dummy).Request().Post(c.ctx)
	if err != nil {
		c.ChangeUserStatus(err)
		return errors.Wrap(err, "msgraph TentativelyAcceptEvent")
	}
	return nil
}
