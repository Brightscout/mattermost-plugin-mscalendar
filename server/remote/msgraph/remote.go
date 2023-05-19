// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"

	msgraph "github.com/yaegashi/msgraph.go/v1.0"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/config"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/serializer"
	"github.com/mattermost/mattermost-plugin-mscalendar/server/utils/bot"
)

const Kind = "msgraph"

type impl struct {
	conf   *config.Config
	logger bot.Logger
}

func init() {
	remote.Makers[Kind] = NewRemote
}

func NewRemote(conf *config.Config, logger bot.Logger) remote.Remote {
	return &impl{
		conf:   conf,
		logger: logger,
	}
}

// MakeClient creates a new client for user-delegated permissions.
func (r *impl) MakeClient(ctx context.Context, token *oauth2.Token, userTokenHelpers *serializer.UserTokenHelpers) remote.Client {
	httpClient := r.NewOAuth2Config().Client(ctx, token)
	c := &client{
		conf:         r.conf,
		ctx:          ctx,
		httpClient:   httpClient,
		Logger:       r.logger,
		rbuilder:     msgraph.NewClient(httpClient),
		tokenHelpers: userTokenHelpers,
	}

	return c
}

// MakeSuperuserClient creates a new client user-delegated permissions with refreshed token.
func (r *impl) MakeUserClient(ctx context.Context, oauthToken *oauth2.Token, mattermostUserID string, userTokenHelpers *serializer.UserTokenHelpers) remote.Client {
	config := r.NewOAuth2Config()

	token, err := userTokenHelpers.RefreshAndStoreToken(oauthToken, config, mattermostUserID)
	if err != nil {
		r.logger.Warnf("Not able to refresh or store the token", "error", err.Error())
		return &client{}
	}

	return r.MakeClient(ctx, token, userTokenHelpers)
}

// MakeSuperuserClient creates a new client used for app-only permissions.
func (r *impl) MakeSuperuserClient(ctx context.Context) (remote.Client, error) {
	httpClient := &http.Client{}
	c := &client{
		conf:       r.conf,
		ctx:        ctx,
		httpClient: httpClient,
		Logger:     r.logger,
		rbuilder:   msgraph.NewClient(httpClient),
	}
	token, err := c.GetSuperuserToken()
	if err != nil {
		return nil, err
	}

	o := &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}
	return r.MakeClient(ctx, o, nil), nil
}

func (r *impl) NewOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     r.conf.OAuth2ClientID,
		ClientSecret: r.conf.OAuth2ClientSecret,
		RedirectURL:  r.conf.PluginURL + config.FullPathOAuth2Redirect,
		Scopes: []string{
			"offline_access",
			"User.Read",
			"Calendars.ReadWrite",
			"Calendars.ReadWrite.Shared",
			"MailboxSettings.Read",
		},
		Endpoint: microsoft.AzureADEndpoint(r.conf.OAuth2Authority),
	}
}
