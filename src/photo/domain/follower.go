package domain

type Follower struct {
	TimeStamp

	Id         string `json:"id,omitempty" gorethink:"id,omitempty"`
	AccountId  string `json:"account_id" gorethink:"account_id"`
	FollowerId string `json:"follower_id" gorethink:"follower_id"`
}
