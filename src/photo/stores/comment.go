package stores

import (
	"errors"
	"photo/api/rethink"
	"photo/domain"
	"photo/utils/refine"
	"photo/utils/text"
	"time"

	r "github.com/dancannon/gorethink"
)

type CommentStore struct {
	re *rethink.Instance
}

func NewCommentStore(re *rethink.Instance) *CommentStore {
	return &CommentStore{
		re: re,
	}
}

func (this *CommentStore) Get(id string) (*domain.Comment, error) {
	var comment *domain.Comment
	err := this.re.One(this.re.Table(kTableComment).Get(id), &comment)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (this *CommentStore) Create(comment *domain.Comment) (*domain.Comment, error) {
	comment.TimeStamp.CreatedTime = time.Now()
	comment.TimeStamp.UpdatedTime = time.Now()
	comment.Tags = text.DeduplicateTag(text.ToTagArray(comment.Text))

	res, err := this.re.RunWrite(this.re.Table(kTableComment).Insert(comment))

	if err != nil {
		return nil, err
	}

	comment.Id = res.GeneratedKeys[0]
	return comment, nil
}

func (this *CommentStore) Update(id string, comment *domain.Comment) (*domain.Comment, error) {
	comment.TimeStamp.UpdatedTime = time.Now()
	comment.Tags = text.DeduplicateTag(text.ToTagArray(comment.Text))

	res, err := this.re.RunWrite(this.re.Table(kTableComment).Get(id).Update(comment))
	if err != nil {
		return nil, err
	}

	if res.Replaced == 0 {
		return nil, errors.New("Update fail")
	}

	return comment, nil
}

func (this *CommentStore) Count(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) (total int, err error) {

	query := refine.RefineCountQuery(this.re, refine.Criteria{
		Table:   kTableComment,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	err = this.re.One(query, &total)

	return
}

func (this *CommentStore) List(index string, keyword interface{}, cond func(r r.Term) r.Term, paging domain.Paging) ([]*domain.Comment, error) {

	query := refine.RefineListQuery(this.re, refine.Criteria{
		Table:   kTableComment,
		Index:   index,
		Keyword: keyword,
		Cond:    cond,
		Paging:  paging,
	})

	var comments []*domain.Comment
	err := this.re.All(query, &comments)

	if err != nil {
		return []*domain.Comment{}, err
	}

	return comments, nil
}

func (this *CommentStore) Delete(id string) (err error) {
	_, err = this.re.Run(this.re.Table(kTableComment).Get(id).Delete())
	return
}

func (this *CommentStore) GetAllKeywords() ([]string, error) {
	var keywords []string
	err := this.re.All(this.re.Table(kTableComment).Field(KIndexCommentByTags).ConcatMap(func(keys r.Term) r.Term {
		return keys
	}).Distinct(), &keywords)

	if err != nil {
		return []string{}, err
	}

	return keywords, err
}
