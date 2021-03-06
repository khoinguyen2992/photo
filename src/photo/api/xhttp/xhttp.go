package xhttp

import (
	"encoding/json"
	"net/http"
	"photo/domain"
	"photo/utils/text"
	"regexp"
	"strconv"
	"strings"

	gorillaContext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type keyContextType int

const keyContext keyContextType = 0

type Context struct {
	Params httprouter.Params
}

func GetContext(r *http.Request) *Context {
	ctx := gorillaContext.Get(r, keyContext).(*Context)
	return ctx
}

func ParseJsonBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&v)
}

func MakeSearchValue(conds []string, keywords []string) []string {
	res := make([]string, 0)
	for i := range conds {
		r, err := regexp.Compile("^" + conds[i])

		if err != nil {
			res := make([]string, 0)
			return res
		}

		for j := range keywords {
			if r.MatchString(keywords[j]) {
				res = append(res, keywords[j])
			}
		}
	}
	return text.Deduplicate(res)
}

type PagingQuery struct {
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Sort     []string `json:"sort"`
}

func NewPagingQuery(r *http.Request) *PagingQuery {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")
	sort := r.URL.Query().Get("sort")

	var p, pS int
	p, err := strconv.Atoi(page)

	if err != nil || p < 0 {
		p = 0
	}

	pS, err = strconv.Atoi(pageSize)

	if err != nil || pS <= 0 {
		pS = 10
	}

	order := make([]string, 0)
	if sort != "" {
		arr := strings.Split(sort, ",")
		for i := range arr {
			order = append(order, arr[i])
		}
	}

	return &PagingQuery{
		Page:     p,
		PageSize: pS,
		Sort:     order,
	}
}

func (p *PagingQuery) AddSort(name string) {
	for _, sort := range p.Sort {
		if sort == name {
			return
		}
	}

	p.Sort = append(p.Sort, name)
}

func (p PagingQuery) Normalize() domain.Paging {
	var pSize = 10
	if p.PageSize > 0 {
		pSize = p.PageSize
	}

	return domain.Paging{
		Start: p.Page * pSize,
		Limit: pSize,
		Sort:  p.Sort,
	}
}
