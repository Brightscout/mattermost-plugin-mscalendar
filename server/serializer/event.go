package serializer

type User struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
	Mail              string `json:"mail,omitempty"`
}

type Event struct {
	Start                      *DateTime            `json:"start,omitempty"`
	Location                   *Location            `json:"location,omitempty"`
	End                        *DateTime            `json:"end,omitempty"`
	Organizer                  *Attendee            `json:"organizer,omitempty"`
	Body                       *ItemBody            `json:"Body,omitempty"`
	ResponseStatus             *EventResponseStatus `json:"responseStatus,omitempty"`
	Importance                 string               `json:"importance,omitempty"`
	ICalUID                    string               `json:"iCalUId,omitempty"`
	Subject                    string               `json:"subject,omitempty"`
	BodyPreview                string               `json:"bodyPreview,omitempty"`
	ShowAs                     string               `json:"showAs,omitempty"`
	Weblink                    string               `json:"weblink,omitempty"`
	ID                         string               `json:"id,omitempty"`
	Attendees                  []*Attendee          `json:"attendees,omitempty"`
	ReminderMinutesBeforeStart int                  `json:"reminderMinutesBeforeStart,omitempty"`
	IsOrganizer                bool                 `json:"isOrganizer,omitempty"`
	IsCancelled                bool                 `json:"isCancelled,omitempty"`
	IsAllDay                   bool                 `json:"isAllDay,omitempty"`
	ResponseRequested          bool                 `json:"responseRequested,omitempty"`
}

type Location struct {
	DisplayName  string       `json:"displayName,omitempty"`
	Address      *Address     `json:"address"`
	Coordinates  *Coordinates `json:"coordinates"`
	LocationType string       `json:"locationType"`
}

type Address struct {
	Street          string `json:"street,omitempty"`
	City            string `json:"city,omitempty"`
	State           string `json:"state,omitempty"`
	CountryOrRegion string `json:"countryOrRegion,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
}

type Attendee struct {
	Status       *EventResponseStatus `json:"status,omitempty"`
	EmailAddress *EmailAddress        `json:"emailAddress,omitempty"`
	Type         string               `json:"type,omitempty"`
}

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}

type EventResponseStatus struct {
	Response string `json:"response,omitempty"`
	Time     string `json:"time,omitempty"`
}

type ItemBody struct {
	Content     string `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}
