package domain

type Profile struct {
	Avatar    string      `json:"avatar" gorethink:"avatar"`
	FirstName string      `json:"first_name" gorethink:"first_name"`
	LastName  string      `json:"last_name" gorethink:"last_name"`
	Followers []*Follower `json:"followers" gorethink:"-"`
}

type Account struct {
	TimeStamp
	Profile

	Id       string `json:"id,omitempty" gorethink:"id,omitempty"`
	Username string `json:"username" gorethink:"username"`
	Salt     string `json:"-" gorethink:"salt"`
	Secret   string `json:"-" gorethink:"secret"`
}

type RegisterAccount struct {
	Username string `json:"username" gorethink:"username"`
	Password string `json:"password" gorethink:"password"`
}

type LoginAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordAccount struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
