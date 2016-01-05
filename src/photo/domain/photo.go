package domain

type Photo struct {
	TimeStamp

	Id        string `json:"id,omitempty" gorethink:"id,omitempty"`
	AccountId string `json:"account_id" gorethink:"account_id"`
	SaveName  string `json:"save_name" gorethink:"save_name"`
	Uri       string `json:"uri" gorethink:"uri"`
	IsPrivate bool   `json:"-" gorethink:"is_private"`
}
