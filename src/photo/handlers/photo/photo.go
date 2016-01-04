package photo

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"photo/api/xhttp"
	"photo/domain"
	"photo/handlers/account"
	"photo/stores"
	"photo/utils/logs"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/dancannon/gorethink"
)

var l = logs.New("handlers/photo")

type PhotoCtrl struct {
	photoStore *stores.PhotoStore
	uploadDir  string
	pathPrefix string
}

func NewPhotoCtrl(photoStore *stores.PhotoStore, uploadDir string, pathPrefix string) *PhotoCtrl {
	return &PhotoCtrl{
		photoStore: photoStore,
		uploadDir:  uploadDir,
		pathPrefix: pathPrefix,
	}
}

func (this *PhotoCtrl) Create(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)
	formFile, formHead, err := r.FormFile("file")

	if err != nil {
		l.Println("Error read form", err)
		xhttp.ResponseBadRequest(w, xhttp.ParseFormFail)
		return
	}
	defer formFile.Close()

	saveName := HashName(formHead.Filename, GenSalt())
	savePath := filepath.Join(this.uploadDir, saveName)

	osFile, err := os.Create(savePath)
	if err != nil {
		l.Println("Error create file", err, savePath)
		xhttp.ResponseBadRequest(w, xhttp.CreateFileFail)
		return
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, formFile)
	if err != nil {
		l.Println("Error save file", err, savePath)
		xhttp.ResponseBadRequest(w, xhttp.CreateFileFail)
		return
	}

	photo, err := this.photoStore.Create(&domain.Photo{
		AccountId: account.Id,
		SaveName:  saveName,
		Uri:       this.pathPrefix + saveName,
		IsPrivate: false,
	})

	if err != nil {
		l.Println("Create", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, photo)
}

func (this *PhotoCtrl) CreateCropImage(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)
	ctx := xhttp.GetContext(r)
	width := ctx.Params.ByName("width")
	height := ctx.Params.ByName("height")
	saveName := ctx.Params.ByName("savename")

	img, err := imaging.Open(path.Join(this.uploadDir, saveName))

	if err != nil {
		l.Println("Crop", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestQuery)
		return
	}

	wd, err := strconv.Atoi(width)

	if err != nil {
		l.Println("Crop", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestQuery)
		return
	}

	he, err := strconv.Atoi(height)

	if err != nil {
		l.Println("Crop", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.InvalidRequestQuery)
		return
	}

	centerCrop := imaging.CropCenter(img, wd, he)

	newName := HashName(width+height+saveName, GenSalt())
	newPath := path.Join(this.uploadDir, newName)

	err = imaging.Save(centerCrop, newPath)

	if err != nil {
		l.Println("Error save file", err, newPath)
		xhttp.ResponseBadRequest(w, xhttp.CreateFileFail)
		return
	}

	photo, err := this.photoStore.Create(&domain.Photo{
		AccountId: account.Id,
		SaveName:  newName,
		Uri:       this.pathPrefix + newName,
		IsPrivate: false,
	})

	if err != nil {
		l.Println("Crop", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, photo)
}

func (this *PhotoCtrl) ListByOwner(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	cond := func(photo gorethink.Term) gorethink.Term {
		return photo
	}

	photos, err := this.photoStore.List(stores.KIndexPhotoByAccountId, account.Id, cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.photoStore.Count(stores.KIndexPhotoByAccountId, account.Id, cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photos == nil {
		photos = []*domain.Photo{}
	}

	xhttp.ResponseList(w, http.StatusOK, photos, pagingReponse, total)
}

func (this *PhotoCtrl) ListByFollowers(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	cond := func(photo gorethink.Term) gorethink.Term {
		origin := photo

		photo = origin.Field(stores.KFilterPhotoByIsPrivate).Eq(false)

		if len(account.Followers) >= 1 {
			follow := origin.Field(stores.KFilterPhotoByAccountId).Eq(account.Followers[0].AccountId)
			for i := range account.Profile.Followers {
				follow = follow.Or(origin.Field(stores.KFilterPhotoByAccountId).Eq(account.Profile.Followers[i].AccountId))
			}
			photo = photo.And(follow)
		}

		return photo
	}

	photos, err := this.photoStore.List(stores.KIndexPhotoByAccountId, "", cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.photoStore.Count(stores.KIndexPhotoByAccountId, "", cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photos == nil {
		photos = []*domain.Photo{}
	}

	xhttp.ResponseList(w, http.StatusOK, photos, pagingReponse, total)
}

func (this *PhotoCtrl) ListByAccount(w http.ResponseWriter, r *http.Request) {
	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	pagingQuery := xhttp.NewPagingQuery(r)
	pagingReponse := pagingQuery.Normalize()

	cond := func(photo gorethink.Term) gorethink.Term {
		photo = photo.Field(stores.KFilterPhotoByIsPrivate).Eq(false)
		return photo
	}

	photos, err := this.photoStore.List(stores.KIndexPhotoByAccountId, id, cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	total, err := this.photoStore.Count(stores.KIndexPhotoByAccountId, id, cond, pagingReponse)

	if err != nil {
		l.Println("List", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photos == nil {
		photos = []*domain.Photo{}
	}

	xhttp.ResponseList(w, http.StatusOK, photos, pagingReponse, total)
}

func (this *PhotoCtrl) UpdatePrivate(w http.ResponseWriter, r *http.Request) {
	account := account.GetAccount(r)

	ctx := xhttp.GetContext(r)
	id := ctx.Params.ByName("id")

	photo, err := this.photoStore.Get(id)

	if err != nil {
		l.Println("UpdatePrivate", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	if photo.AccountId != account.Id {
		l.Println("UpdatePrivate", r.URL, errors.New("No permission"))
		xhttp.ResponseBadRequest(w, xhttp.NoPermission)
		return
	}

	photo.IsPrivate = !photo.IsPrivate

	photo, err = this.photoStore.Update(id, photo)

	if err != nil {
		l.Println("UpdatePrivate", r.URL, err)
		xhttp.ResponseBadRequest(w, xhttp.DatabaseError)
		return
	}

	xhttp.ResponseJson(w, http.StatusOK, photo)
}

func GenSalt() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(rand.Int63()))))
}

func HashName(filename string, salt string) string {
	sum := sha256.Sum256([]byte(filename + salt))
	return hex.EncodeToString(sum[:]) + "_" + filename
}
