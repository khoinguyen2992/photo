package stores

import (
	"time"

	"photo/api/rethink"
	"photo/domain"
	"photo/utils/refine"

	r "github.com/dancannon/gorethink"
)

type FollowerStore struct {
	re *rethink.Instance
}

func NewFollowerStore(re *rethink.Instance) *FollowerStore {
	return &FollowerStore{
		re: re,
	}
}
func (this *FollowerStore) Get(id string) (*domain.Follower, error) {
	var follower *domain.Follower
	err := this.re.One(this.re.Table(kTableFollower).Get(id), &follower)

	if err != nil {
		return nil, err
	}

	return follower, nil
}

func (this *FollowerStore) IsExisted(a, f string) (*domain.Account, error) {
	var follower *domain.Account

	query := make(map[string]interface{})
	query[KIndexFollowerByFollowerId] = f
	err := this.re.One(this.re.Table(kTableFollower).GetAllByIndex(KIndexFollowerByAccountId, a).Filter(query), &follower)

	if err != nil {
		return nil, err
	}

	return follower, nil
}

func (this *FollowerStore) Count(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (total int, err error) {

	query := refine.RefineCountQuery(this.re, refine.Criteria{
		Table:   kTableFollower,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.One(query, &total)

	return
}

func (this *FollowerStore) List(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (followers []*domain.Follower, err error) {

	query := refine.RefineListQuery(this.re, refine.Criteria{
		Table:   kTableFollower,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.All(query, &followers)

	return
}

func (this *FollowerStore) ListAll(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (followers []*domain.Follower, err error) {

	query := refine.RefineListQueryAll(this.re, refine.Criteria{
		Table:   kTableFollower,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.All(query, &followers)

	return
}

func (this *FollowerStore) Create(follower *domain.Follower) (*domain.Follower, error) {
	follower.TimeStamp.CreatedTime = time.Now()
	follower.TimeStamp.UpdatedTime = time.Now()

	res, err := this.re.RunWrite(this.re.Table(kTableFollower).Insert(follower))

	if err != nil {
		return nil, err
	}

	follower.Id = res.GeneratedKeys[0]
	return follower, nil
}

func (this *FollowerStore) Delete(id string) (err error) {
	_, err = this.re.Run(this.re.Table(kTableFollower).Get(id).Delete())
	return
}
