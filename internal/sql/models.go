// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Journal struct {
	Issn  string `json:"issn"`
	Title string `json:"title"`
	Short string `json:"short"`
}

type Paper struct {
	ID       int64       `json:"id"`
	Doi      string      `json:"doi"`
	Url      pgtype.Text `json:"url"`
	Issn     string      `json:"issn"`
	Title    string      `json:"title"`
	Crossref []byte      `json:"crossref"`
	Body     pgtype.Text `json:"body"`
}
