package domain

import "time"

type TimeStamp struct {
	CreatedTime time.Time `json:"created_time,omitempty" gorethink:"created_time,omitempty"`
	UpdatedTime time.Time `json:"updated_time,omitempty" gorethink:"updated_time,omitempty"`
}

type Paging struct {
	Start int      `json:"start"`
	Limit int      `json:"limit"`
	Sort  []string `json:"sort"`
	Total int      `json:"total"`
}
