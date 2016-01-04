package stores

import (
	"errors"
	"photo/api/rethink"
	"photo/domain"
	"photo/utils/refine"
	"time"

	r "github.com/dancannon/gorethink"
)

type PhotoStore struct {
	re *rethink.Instance
}

func NewPhotoStore(re *rethink.Instance) *PhotoStore {
	return &PhotoStore{
		re: re,
	}
}

func (this *PhotoStore) Get(id string) (*domain.Photo, error) {
	var photo *domain.Photo
	err := this.re.One(this.re.Table(kTablePhoto).Get(id), &photo)

	if err != nil {
		return nil, err
	}

	return photo, nil
}

func (this *PhotoStore) Count(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (total int, err error) {

	query := refine.RefineCountQuery(this.re, refine.Criteria{
		Table:   kTablePhoto,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.One(query, &total)

	return
}

func (this *PhotoStore) List(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) ([]*domain.Photo, error) {

	query := refine.RefineListQuery(this.re, refine.Criteria{
		Table:   kTablePhoto,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	var photos []*domain.Photo
	err := this.re.All(query, &photos)

	if err != nil {
		return []*domain.Photo{}, err
	}

	return photos, nil
}

func (this *PhotoStore) Create(photo *domain.Photo) (*domain.Photo, error) {
	photo.TimeStamp.CreatedTime = time.Now()
	photo.TimeStamp.UpdatedTime = time.Now()

	res, err := this.re.RunWrite(this.re.Table(kTablePhoto).Insert(photo))

	if err != nil {
		return nil, err
	}

	photo.Id = res.GeneratedKeys[0]
	return photo, nil
}

func (this *PhotoStore) Update(id string, photo *domain.Photo) (*domain.Photo, error) {
	photo.TimeStamp.UpdatedTime = time.Now()

	res, err := this.re.RunWrite(this.re.Table(kTablePhoto).Get(id).Update(photo))
	if err != nil {
		return nil, err
	}

	if res.Replaced == 0 {
		return nil, errors.New("Update fail")
	}

	return photo, nil
}

func (this *PhotoStore) Delete(id string) (err error) {
	_, err = this.re.Run(this.re.Table(kTablePhoto).Get(id).Delete())
	return
}
