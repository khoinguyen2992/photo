package stores

import (
	"errors"
	"time"

	"photo/api/rethink"
	"photo/domain"
	"photo/utils/refine"

	r "github.com/dancannon/gorethink"
)

type AccountStore struct {
	re *rethink.Instance
}

func NewAccountStore(re *rethink.Instance) *AccountStore {
	return &AccountStore{
		re: re,
	}
}

func (this *AccountStore) Get(id string) (*domain.Account, error) {
	var account *domain.Account

	err := this.re.One(this.re.Table(kTableAccount).Get(id), &account)
	if err != nil {
		return nil, err
	}

	var followers []*domain.Follower
	err = this.re.All(this.re.Table(kTableFollower).GetAllByIndex(KIndexFollowerByFollowerId, id), &followers)

	if err != nil {
		account.Profile.Followers = []*domain.Follower{}
	} else {
		account.Profile.Followers = followers
	}

	return account, nil
}

func (this *AccountStore) Count(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (total int, err error) {

	query := refine.RefineCountQuery(this.re, refine.Criteria{
		Table:   kTableAccount,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.One(query, &total)

	return
}

func (this *AccountStore) List(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) ([]*domain.Account, error) {

	query := refine.RefineListQuery(this.re, refine.Criteria{
		Table:   kTableAccount,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	var accounts []*domain.Account
	err := this.re.All(query, &accounts)

	if err != nil {
		return []*domain.Account{}, err
	}

	return accounts, nil
}

func (this *AccountStore) IsExisted(username string) (*domain.Account, error) {
	var account *domain.Account

	err := this.re.One(this.re.Table(kTableAccount).GetAllByIndex(KIndexAccountByUsername, username), &account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (this *AccountStore) Create(account *domain.Account) (*domain.Account, error) {
	account.TimeStamp.CreatedTime = time.Now()
	account.TimeStamp.UpdatedTime = time.Now()

	res, err := this.re.RunWrite(this.re.Table(kTableAccount).Insert(account))

	if err != nil {
		return nil, err
	}

	account.Id = res.GeneratedKeys[0]
	account.Profile.Followers = []*domain.Follower{}
	return account, nil
}

func (this *AccountStore) Update(id string, account *domain.Account) (*domain.Account, error) {
	account.TimeStamp.UpdatedTime = time.Now()

	res, err := this.re.RunWrite(this.re.Table(kTableAccount).Get(id).Update(account))
	if err != nil {
		return nil, err
	}

	if res.Replaced == 0 {
		return nil, errors.New("Update fail")
	}

	return account, nil
}

func (this *AccountStore) Delete(id string) (err error) {
	_, err = this.re.Run(this.re.Table(kTableAccount).Get(id).Delete())
	return
}
