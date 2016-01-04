package comment

import (
	"errors"
	"net/http"
	"photo/api/xhttp"
	"photo/domain"
	"photo/handlers/account"
	"photo/stores"
	"photo/utils/logs"
	"photo/utils/text"
	"strings"

	"github.com/dancannon/gorethink"
)

var l = logs.New("handlers/comment")

type CommentCtrl struct {
	commentStore *stores.CommentStore
	photoStore   *stores.PhotoStore
}

func NewCommentCtrl(commentStore *stores.CommentStore, photoStore *stores.PhotoStore) *CommentCtrl {
	return &CommentCtrl{
		commentStore: commentStore,
		photoStore:   photoStore,
	}
}

func (this *CommentCtrl) ListByTags(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	keywords, err := this.commentStore.GetAllKeywords()
	if err != nil {
		l.Println("ListByTag", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	word := ""
	searchValue := make([]string, 0)

	if tags != "" {
		cond := text.ToArray(text.Normalize(tags), ",")

		searchValue = xhttp.MakeSearchValue(cond, keywords)

		if len(searchValue) == 0 {
			l.Println("List", r.URL, "No keywords matched")
			xhttp.ResponseList(w, http.StatusOK, []interface{}{}, pagingReponse, 0)
			return
		}

		if len(searchValue) > 0 && len(cond) > 1 {
			word = strings.Split(searchValue[0], ",")[0]
		}
	}

	cond := func(comment gorethink.Term) gorethink.Term {
		origin := comment
		if tags != "" {
			values := text.StringToInterfaceArray(text.ToArray(searchValue[0], ","))
			commentTag := origin.Field(stores.KIndexCommentByTags).Contains(values...)
			for i := 0; i < len(searchValue); i++ {
				values := text.StringToInterfaceArray(text.ToArray(searchValue[i], ","))
				commentTag = commentTag.Or(origin.Field(stores.KIndexCommentByTags).Contains(values...))
			}

			comment = comment.And(commentTag)
		}
		return comment
	}

	comments, err := this.commentStore.List(stores.KIndexCommentByPhotoId, word, cond, pagingReponse)

	if err != nil {
		l.Println("ListByTags", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if comments == nil {
		comments = []*domain.Comment{}
	}

	photoIds := make([]string, 0)
	for i := range comments {
		photoIds = append(photoIds, comments[i].PhotoId)
	}

	photoIds = text.Deduplicate(photoIds)

	cond = func(photo gorethink.Term) gorethink.Term {
		origin := photo
		photo = origin.Field(stores.KFilterPhotoByIsPrivate).Eq(false)

		if len(photoIds) >= 1 {
			ph := origin.Field(stores.KIndexPhotoById).Eq(photoIds[0])
			for i := range photoIds {
				ph = ph.Or(origin.Field(stores.KIndexPhotoById).Eq(photoIds[i]))
			}
			photo = photo.And(ph)
		}

		return photo
	}

	photos, err := this.photoStore.List(stores.KIndexPhotoById, "", cond, pagingReponse)

	if err != nil {
		l.Println("ListByTags", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.photoStore.Count(stores.KIndexPhotoById, "", cond, pagingReponse)

	if err != nil {
		l.Println("ListByTags", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseList(w, http.StatusOK, photos, pagingReponse, total)
}

func (this *CommentCtrl) ListByPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	photo, err := this.photoStore.Get(id)

	if err != nil {
		l.Println("ListByPhoto", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photo.IsPrivate {
		l.Println("ListByPhoto", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	cond := func(comment gorethink.Term) gorethink.Term {
		return comment
	}

	comments, err := this.commentStore.List(stores.KIndexCommentByPhotoId, photo.Id, cond, pagingReponse)

	if err != nil {
		l.Println("ListByPhoto", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.commentStore.Count(stores.KIndexCommentByPhotoId, photo.Id, cond, pagingReponse)

	if err != nil {
		l.Println("ListByPhoto", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if comments == nil {
		comments = []*domain.Comment{}
	}

	xhttp.ResponseList(w, http.StatusOK, comments, pagingReponse, total)
}

func (this *CommentCtrl) ListByNotification(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	cond := func(comment gorethink.Term) gorethink.Term {
		comment = comment.Field(stores.KFilterCommentByIsKnown).Eq(false)
		return comment
	}

	comments, err := this.commentStore.List(stores.KIndexCommentByNotificationId, account.Id, cond, pagingReponse)

	if err != nil {
		l.Println("ListByNoti", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.commentStore.Count(stores.KIndexCommentByPhotoId, account.Id, cond, pagingReponse)

	if err != nil {
		l.Println("ListByNoti", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if comments == nil {
		comments = []*domain.Comment{}
	}

	xhttp.ResponseList(w, http.StatusOK, comments, pagingReponse, total)
}

func (this *CommentCtrl) Create(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	var comment *domain.Comment
	err := xhttp.ParseJsonBody(r, &comment)

	if err != nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestBody)
		return
	}

	photo, err := this.photoStore.Get(comment.PhotoId)

	if err != nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photo.IsPrivate {
		l.Println("Create", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	comment.AccountId = account.Id
	comment.NotificationId = photo.AccountId
	comment.IsKnown = false

	comment, err = this.commentStore.Create(comment)

	if err != nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, comment)
}

func (this *CommentCtrl) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	comment, err := this.commentStore.Get(id)
	if err != nil {
		l.Println("UpdateNoti", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	photo, err := this.photoStore.Get(comment.PhotoId)

	if err != nil {
		l.Println("UpdateNoti", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photo.AccountId != account.Id {
		l.Println("UpdateNoti", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	comment.IsKnown = true

	comment, err = this.commentStore.Update(id, comment)

	if err != nil {
		l.Println("UpdateNoti", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, comment)
}

func (this *CommentCtrl) Delete(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	comment, err := this.commentStore.Get(id)
	if err != nil {
		l.Println("Delete", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if comment.AccountId != account.Id || comment.NotificationId != account.Id {
		l.Println("Delete", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	err = this.commentStore.Delete(id)

	if err != nil {
		l.Println("Delete", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseOk(w)
}
