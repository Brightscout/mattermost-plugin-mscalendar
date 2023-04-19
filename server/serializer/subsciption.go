package serializer

type Subscription struct {
	ID                 string `json:"id"`
	Resource           string `json:"resource,omitempty"`
	ApplicationID      string `json:"applicationId,omitempty"`
	ChangeType         string `json:"changeType,omitempty"`
	ClientState        string `json:"clientState,omitempty"`
	NotificationURL    string `json:"notificationUrl,omitempty"`
	ExpirationDateTime string `json:"expirationDateTime,omitempty"`
	CreatorID          string `json:"creatorId,omitempty"`
}
