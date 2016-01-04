package account

import (
	"net/http"
	"photo/api/xhttp"
	"photo/domain"
	"photo/services/auth"
	"photo/stores"
	"photo/utils/logs"
	"strings"

	"github.com/asaskevich/govalidator"

	"github.com/dancannon/gorethink"
	gorillaContext "github.com/gorilla/context"
)

var l = logs.New("controller/account")

type AccountCtrl struct {
	accountStore *stores.AccountStore
	accountAuth  *auth.AccountAuth
}

func NewAccountCtrl(accountStore *stores.AccountStore, accountAuth *auth.AccountAuth) *AccountCtrl {
	return &AccountCtrl{
		accountStore: accountStore,
		accountAuth:  accountAuth,
	}
}

func (this *AccountCtrl) Logout(w http.ResponseWriter, r *http.Request) {
	session := this.accountAuth.CurrentSession(r)
	this.accountAuth.SessionClear(session)
	this.accountAuth.SessionSave(session, w, r)
	xhttp.ResponseOk(w)
}

func (this *AccountCtrl) Login(w http.ResponseWriter, r *http.Request) {
	currentSession := this.accountAuth.CurrentSession(r)
	account_id := this.accountAuth.SessionGet(currentSession)
	account, err := this.accountStore.Get(account_id)

	if err == nil {
		l.Println("Login", r.URL, "Already session")
		this.accountAuth.SessionSet(currentSession, string(account.Id))
		this.accountAuth.SessionSave(currentSession, w, r)
		xhttp.ResponseJson(w, http.StatusOK, account)
		return
	}

	l.Println("Login", r.URL, "No session")
	var loginAccount *domain.LoginAccount
	err = xhttp.ParseJsonBody(r, &loginAccount)

	if err != nil {
		l.Println("Login", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestBody)
		return
	}

	loginAccount.Username = strings.TrimSpace(loginAccount.Username)
	loginAccount.Password = strings.TrimSpace(loginAccount.Password)
	if loginAccount.Username == "" || loginAccount.Password == "" || !govalidator.IsEmail(loginAccount.Username) {
		l.Println("Login", r.URL, "Missing fields")
		xhttp.ResponseBadRequest(w, xhttp.EmptyRequiredFields)
		return
	}

	account, err = this.accountStore.IsExisted(loginAccount.Username)

	if err != nil {
		l.Println("Login", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.UsernameNotExisted)
		return
	}

	newSecret := auth.HashPassword(account.Salt, loginAccount.Password, loginAccount.Username)

	if newSecret != account.Secret {
		l.Println("Login", r.URL, "Invalid password")
		xhttp.ResponseBadRequest(w, xhttp.InvalidPassword)
		return
	}

	this.accountAuth.SessionSet(currentSession, string(account.Id))
	this.accountAuth.SessionSave(currentSession, w, r)
	xhttp.ResponseJson(w, http.StatusOK, account)
	return
}

func (this *AccountCtrl) Register(w http.ResponseWriter, r *http.Request) {
	var registerAccount *domain.RegisterAccount
	err := xhttp.ParseJsonBody(r, &registerAccount)

	if err != nil {
		l.Println("Register", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestBody)
		return
	}

	registerAccount.Username = strings.TrimSpace(registerAccount.Username)
	registerAccount.Password = strings.TrimSpace(registerAccount.Password)

	if registerAccount.Username == "" || registerAccount.Password == "" || !govalidator.IsEmail(registerAccount.Username) {
		l.Println("Register", r.URL, "Missing fields")
		xhttp.ResponseBadRequest(w, xhttp.EmptyRequiredFields)
		return
	}

	existedAccount, err := this.accountStore.IsExisted(registerAccount.Username)

	if err == nil {
		l.Println("Register", r.URL, "Duplicated email", existedAccount.Username)
		xhttp.ResponseBadRequest(w, xhttp.DuplicatedUsername)
		return
	}

	newSalt := auth.GenSalt()
	newSecret := auth.HashPassword(newSalt, registerAccount.Password, registerAccount.Username)

	newAccount := &domain.Account{
		Username: registerAccount.Username,
		Salt:     newSalt,
		Secret:   newSecret,
	}

	newAccount, err = this.accountStore.Create(newAccount)

	if err != nil {
		l.Println("Register", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, newAccount)
	return

}

func (this *AccountCtrl) Me(w http.ResponseWriter, r *http.Request) {
	account := GetAccount(r)
	xhttp.ResponseJson(w, http.StatusOK, account)
}

func (this *AccountCtrl) List(w http.ResponseWriter, r *http.Request) {

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	cond := func(account gorethink.Term) gorethink.Term {
		return account
	}

	accounts, err := this.accountStore.List(stores.KIndexAccountByUsername, "", cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.accountStore.Count(stores.KIndexAccountByUsername, "", cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if accounts == nil {
		accounts = []*domain.Account{}
	}

	xhttp.ResponseList(w, http.StatusOK, accounts, pagingReponse, total)
}

func (this *AccountCtrl) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	account := GetAccount(r)

	var passwordAccount *domain.PasswordAccount
	err := xhttp.ParseJsonBody(r, &passwordAccount)

	if err != nil {
		l.Println("UpdatePassword", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestBody)
		return
	}

	oldSecret := auth.HashPassword(account.Salt, passwordAccount.OldPassword, account.Username)

	if oldSecret != account.Secret {
		l.Println("UpdatePassword", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidPassword)
		return
	}

	newSalt := auth.GenSalt()
	newSecret := auth.HashPassword(newSalt, passwordAccount.NewPassword, account.Username)
	account.Salt = newSalt
	account.Secret = newSecret

	account, err = this.accountStore.Update(account.Id, account)

	if err != nil {
		l.Println("Update", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, account)
}

func (this *AccountCtrl) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	account := GetAccount(r)

	var profile *domain.Profile
	err := xhttp.ParseJsonBody(r, &profile)

	if err != nil {
		l.Println("UpdateProfile", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestBody)
		return
	}

	account.Profile.FirstName = profile.FirstName
	account.Profile.LastName = profile.LastName
	account.Profile.Avatar = profile.Avatar

	newAccount, err := this.accountStore.Update(account.Id, account)

	if err != nil {
		l.Println("UpdateProfile", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, newAccount)
	return
}

func GetAccount(r *http.Request) *domain.Account {
	account, ok := gorillaContext.Get(r, "account").(*domain.Account)

	if ok {
		return account
	}

	return nil
}
