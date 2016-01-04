package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	KeyUserId = "userid"
)

type AccountAuth struct {
	SessionName  string
	SessionStore *sessions.CookieStore
}

func NewAccountAuth(sessionName string, cookieStore string) *AccountAuth {
	return &AccountAuth{
		SessionName:  sessionName,
		SessionStore: sessions.NewCookieStore([]byte(cookieStore)),
	}
}

func (this *AccountAuth) CurrentSession(r *http.Request) *sessions.Session {
	s, _ := this.SessionStore.Get(r, this.SessionName)
	return s
}

func (this *AccountAuth) SessionGet(session *sessions.Session) string {
	if v, ok := session.Values[KeyUserId]; ok {
		return v.(string)
	}

	return ""
}

func (this *AccountAuth) SessionClear(session *sessions.Session) {
	delete(session.Values, KeyUserId)
}

func (this *AccountAuth) SessionSet(session *sessions.Session, value string) {
	session.Values[KeyUserId] = value
}

func (this *AccountAuth) SessionSave(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Save(r, w)
}

func GenSalt() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(rand.Int63()))))
}

func HashPassword(salt, password, username string) string {
	sum := sha256.Sum256([]byte(username + password + salt))
	return hex.EncodeToString(sum[:])
}
