package xhttp

const (
	//common
	InvalidRequestBody  = "invalid request body"
	InvalidRequestQuery = "invalid request query"
	EmptyRequiredFields = "empty required fields"
	DatabaseError       = "database error"
	NoPermission        = "no permission"
	NoRecordFound       = "no record found"

	//auth
	InvalidPassword    = "invalid password"
	UsernameNotExisted = "username not existed"
	DuplicatedUsername = "duplicated username"
	InvalidAccount     = "invalid account"

	//form
	ParseFormFail  = "parse form fail"
	CreateFileFail = "create file fail"
)
