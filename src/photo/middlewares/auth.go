package middlewares

import (
	"net/http"
	"photo/api/xhttp"
	"photo/services/auth"
	"photo/stores"
	"photo/utils/logs"

	gorillaContext "github.com/gorilla/context"
)

var l = logs.New("middlewares/auth")

type Auth struct {
	accountAuth  *auth.AccountAuth
	accountStore *stores.AccountStore
}

func NewAuth(accountStore *stores.AccountStore, accountAuth *auth.AccountAuth) func(http.Handler) http.Handler {
	a := Auth{
		accountAuth:  accountAuth,
		accountStore: accountStore,
	}
	return a.factory
}

func (this Auth) factory(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			currentSession := this.accountAuth.CurrentSession(r)
			account_id := this.accountAuth.SessionGet(currentSession)
			account, err := this.accountStore.Get(account_id)

			if err != nil {
				l.Println(err)
				xhttp.ResponseForbidden(w, "You are not allowed to access this page.")
				return
			}

			gorillaContext.Set(r, "account", account)
			next.ServeHTTP(w, r)
		})
}
