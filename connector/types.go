package connector

import (
	"database/sql"
	"time"
)

type Connector interface {
	Query(db *sql.DB) (*Response, error)
	GetDB(*string, *int, *string, *string, *string, *string) *sql.DB
}

type QueryResponse struct {
	Name         string `json:"name"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Rows         []Row  `json:"rows,omitempty"`
	Duration     int64  `json:"duration"`
	Timestamp    int64  `json:"timestamp"`
}

type Response struct {
	Time         time.Time        `json:"time"`
	Error        bool             `json:"error"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
	Queries      []*QueryResponse `json:"queries,omitempty"`
}

func NewResponse() *Response {
	r := new(Response)
	r.Time = time.Now().UTC()
	return r
}

type Row struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
