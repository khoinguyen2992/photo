package follower

import (
	"errors"
	"net/http"
	"photo/api/xhttp"
	"photo/domain"
	"photo/handlers/account"
	"photo/stores"
	"photo/utils/logs"
	"photo/utils/sendmail"
	"strings"
)

var l = logs.New("handlers/follower")

type FollowerCtrl struct {
	followerStore *stores.FollowerStore
	accountStore  *stores.AccountStore
	sendMail      *sendmail.SendMail
}

func NewFollowerCtrl(followerStore *stores.FollowerStore, accountStore *stores.AccountStore, sendMail *sendmail.SendMail) *FollowerCtrl {
	return &FollowerCtrl{
		followerStore: followerStore,
		accountStore:  accountStore,
		sendMail:      sendMail,
	}
}

func (this *FollowerCtrl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	account := account.GetAccount(r)

	if id == account.Id {
		l.Println("Create", r.URL, "Invalid query")
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestQuery)
		return
	}

	noti, err := this.accountStore.Get(id)

	if err != nil {
		l.Println("Create", r.URL, "Invalid query")
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestQuery)
		return
	}

	_, err = this.followerStore.IsExisted(id, account.Id)
	if err == nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	follower, err := this.followerStore.Create(&domain.Follower{
		AccountId:  id,
		FollowerId: account.Id,
	})

	if err != nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	err = this.sendMail.Send(noti.Username, strings.Join([]string{noti.Profile.FirstName, noti.Profile.LastName}, " "))

	if err != nil {
		l.Println(err)
	}

	xhttp.ResponseJson(w, http.StatusOK, follower)
	return
}

func (this *FollowerCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	account := account.GetAccount(r)

	follower, err := this.followerStore.Get(id)
	if err != nil {
		l.Println("Delete", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if follower.AccountId != account.Id && follower.FollowerId != account.Id {
		l.Println("Delete", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	err = this.followerStore.Delete(id)

	if err != nil {
		l.Println("Delete", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseOk(w)
}
