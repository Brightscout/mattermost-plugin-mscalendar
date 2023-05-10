// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package msgraph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-mscalendar/server/remote"
)

type calendarViewResponse struct {
	Error *remote.APIError `json:"error,omitempty"`
	Value []*remote.Event  `json:"value,omitempty"`
}

type calendarViewSingleResponse struct {
	Headers map[string]string    `json:"headers"`
	ID      string               `json:"id"`
	Body    calendarViewResponse `json:"body"`
	Value   []*remote.Event      `json:"value,omitempty"`
	Status  int                  `json:"status"`
}

type calendarViewBatchResponse struct {
	Responses []*calendarViewSingleResponse `json:"responses"`
}

func (c *client) GetDefaultCalendarView(remoteUserID string, start, end time.Time) ([]*remote.Event, error) {
	paramStr := getQueryParamStringForCalendarView(start, end)

	res := &calendarViewResponse{}
	err := c.rbuilder.Users().ID(remoteUserID).CalendarView().Request().JSONRequest(
		c.ctx, http.MethodGet, paramStr, nil, res)
	if err != nil {
		return nil, errors.Wrap(err, "msgraph GetDefaultCalendarView")
	}

	return res.Value, nil
}

func (c *client) getCalenderView(reqs *singleRequest) (*calendarViewSingleResponse, error) {
	req, err := http.NewRequest(reqs.Method, reqs.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", reqs.AccessToken.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out *calendarViewSingleResponse
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	out.ID = reqs.ID

	return out, err
}

func (c *client) DoBatchViewCalendarRequests(allParams []*remote.ViewCalendarParams) ([]*remote.ViewCalendarResponse, error) {
	requests := []*singleRequest{}
	for _, params := range allParams {
		paramStr := getQueryParamStringForCalendarView(params.StartTime, params.EndTime)
		url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/calendarview%s", paramStr)
		req := &singleRequest{
			ID:          params.RemoteUserID,
			URL:         url,
			Method:      http.MethodGet,
			Headers:     map[string]string{},
			AccessToken: params.AccessToken,
		}

		requests = append(requests, req)
	}

	// batchRequests := prepareBatchRequests(requests)
	// var batchResponses []*calendarViewBatchResponse
	// for _, req := range batchRequests {
	// 	batchRes := &calendarViewBatchResponse{}
	// 	err := c.batchRequest(req, batchRes)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "msgraph ViewCalendar batch request")
	// 	}

	// 	batchResponses = append(batchResponses, batchRes)
	// }

	// result := []*remote.ViewCalendarResponse{}
	// for _, res := range batchResponses {
	// 	for _, res := range batchRes.Responses {
	// 		viewCalRes := &remote.ViewCalendarResponse{
	// 			RemoteUserID: res.ID,
	// 			Events:       res.Body.Value,
	// 			Error:        res.Body.Error,
	// 		}
	// 		result = append(result, viewCalRes)
	// 	}
	// }

	var responses []*calendarViewSingleResponse
	for _, req := range requests {
		resp, err := c.getCalenderView(req)
		if err != nil {
			return nil, errors.Wrap(err, "msgraph ViewCalendar batch request")
		}

		responses = append(responses, resp)
	}

	result := []*remote.ViewCalendarResponse{}
	for _, res := range responses {
		viewCalRes := &remote.ViewCalendarResponse{
			RemoteUserID: res.ID,
			Events:       res.Value,
		}
		result = append(result, viewCalRes)
	}

	return result, nil
}

func getCalendarViewURL(params *remote.ViewCalendarParams) string {
	paramStr := getQueryParamStringForCalendarView(params.StartTime, params.EndTime)
	return "/Users/" + params.RemoteUserID + "/calendarView" + paramStr
}

func getQueryParamStringForCalendarView(start, end time.Time) string {
	q := url.Values{}
	q.Add("startDateTime", start.Format(time.RFC3339))
	q.Add("endDateTime", end.Format(time.RFC3339))
	q.Add("$top", "20")
	return "?" + q.Encode()
}
