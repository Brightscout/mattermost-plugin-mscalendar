// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package remote

type WorkingHours struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TimeZone  struct {
		Name string `json:"name"`
	}
	DaysOfWeek []string `json:"daysOfWeek"`
}

type MailboxSettings struct {
	TimeZone     string       `json:"timeZone"`
	WorkingHours WorkingHours `json:"workingHours"`
}
