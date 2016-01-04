package xhttp

import (
	"encoding/json"
	"net/http"
	"photo/domain"
)

type PagingResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ResponseJson(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)
	w.Write(data)
}

func ResponseForbidden(w http.ResponseWriter, msg string) {
	ResponseJson(w, http.StatusForbidden, ErrorResponse{
		Code:    "forbidden",
		Message: msg,
	})
}

func ResponseBadRequest(w http.ResponseWriter, msg string) {
	ResponseJson(w, http.StatusBadRequest, ErrorResponse{
		Code:    "bad-request",
		Message: msg,
	})
}

func ResponseOk(w http.ResponseWriter) {
	ResponseJson(w, http.StatusOK, ErrorResponse{
		Code:    "ok",
		Message: "ok",
	})
}

func MakePaging(paging domain.Paging) PagingResponse {
	var p = 0
	if paging.Limit > 0 {
		p = int(paging.Start / paging.Limit)
	}
	if p < 0 {
		p = 0
	}

	return PagingResponse{
		Total:    paging.Total,
		Page:     p,
		PageSize: paging.Limit,
	}
}

func ResponseList(w http.ResponseWriter, httpCode int, data interface{}, paging domain.Paging, total int) {
	type Res struct {
		Items  interface{}    `json:"items"`
		Paging PagingResponse `json:"paging"`
	}

	paging.Total = total
	ResponseJson(w, httpCode, &Res{
		Items:  data,
		Paging: MakePaging(paging),
	})
}
