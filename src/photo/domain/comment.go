package domain

type Comment struct {
	TimeStamp

	Id             string   `json:"id,omitempty" gorethink:"id,omitempty"`
	PhotoId        string   `json:"photo_id" gorethink:"photo_id"`
	AccountId      string   `json:"-" gorethink:"account_id"`
	NotificationId string   `json:"-" gorethink:"notification_id"`
	Tags           []string `json:"-" gorethink:"tags"`
	Text           string   `json:"text" gorethink:"text"`
	IsKnown        bool     `json:"-" gorethink:"is_known"`
}
