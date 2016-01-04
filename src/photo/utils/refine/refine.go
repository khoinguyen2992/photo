package refine

import (
	"photo/api/rethink"
	"photo/domain"

	r "github.com/dancannon/gorethink"
)

type Criteria struct {
	Table   string
	Index   string
	Keyword interface{}
	Cond    func(r.Term) r.Term
	Paging  domain.Paging
}

func RefineSort(sort []string) []interface{} {
	if len(sort) == 0 {
		sort = append(sort, "id")
	}

	res := make([]interface{}, len(sort))
	for i := range sort {
		if string(sort[i][0]) == "-" {
			res[i] = r.Desc(sort[i][1:])
		} else {
			res[i] = sort[i]
		}
	}

	return res
}

func RefineListQuery(re *rethink.Instance, criteria Criteria) r.Term {
	sort := RefineSort(criteria.Paging.Sort)

	query := re.Table(criteria.Table)

	if criteria.Keyword != "" {
		query = query.GetAllByIndex(criteria.Index, criteria.Keyword)
	}

	query = query.Filter(criteria.Cond).OrderBy(sort...).Skip(criteria.Paging.Start).Limit(criteria.Paging.Limit)

	return query
}

func RefineListQueryAll(re *rethink.Instance, criteria Criteria) r.Term {
	sort := RefineSort(criteria.Paging.Sort)

	query := re.Table(criteria.Table)

	if criteria.Keyword != "" {
		query = query.GetAllByIndex(criteria.Index, criteria.Keyword)
	}

	query = query.Filter(criteria.Cond).OrderBy(sort...)

	return query
}

func RefineCountQuery(re *rethink.Instance, criteria Criteria) r.Term {

	sort := RefineSort(criteria.Paging.Sort)

	query := re.Table(criteria.Table)

	if criteria.Keyword != "" {
		query = query.GetAllByIndex(criteria.Index, criteria.Keyword)
	}

	query = query.Filter(criteria.Cond).OrderBy(sort...).Count()

	return query
}
