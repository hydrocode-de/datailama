// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: paper.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getPaperStatistics = `-- name: GetPaperStatistics :many
SELECT j.title,
    paper.issn,
    ((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0] AS year,
    count(*) AS count
   FROM paper
     JOIN journals j ON paper.issn = j.issn
  GROUP BY paper.issn, j.title, (((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0])
  ORDER BY (((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0])
`

type GetPaperStatisticsRow struct {
	Title string      `json:"title"`
	Issn  string      `json:"issn"`
	Year  interface{} `json:"year"`
	Count int64       `json:"count"`
}

func (q *Queries) GetPaperStatistics(ctx context.Context) ([]GetPaperStatisticsRow, error) {
	rows, err := q.db.Query(ctx, getPaperStatistics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPaperStatisticsRow
	for rows.Next() {
		var i GetPaperStatisticsRow
		if err := rows.Scan(
			&i.Title,
			&i.Issn,
			&i.Year,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchPaperByTitle = `-- name: SearchPaperByTitle :many
SELECT 
  paper.id,
  paper.title,
  paper.doi,
  paper.url,
  journals.title as journal,
  paper.crossref->'author'->0->>'family' || ', ' || (paper.crossref->'author'->0->>'given')::text as author,
  date(paper.crossref->'published'->'date-parts'->>0) as published,
  crossref->>'is-referenced-by-count' as citations,
  ((crossref->>'is-referenced-by-count')::double precision / (date_part('year', now()) - (paper.crossref->'published'->'date-parts'->0->>0)::double precision + 0.1))::double precision as "citations_year"
FROM paper
JOIN journals ON journals.issn=paper.issn
WHERE paper.title ILIKE '%' || $2::text || '%'
AND paper.crossref->'author'->0->>'family' ILIKE '%' || $3::text || '%'
ORDER BY 
  CASE $4::text
    WHEN 'citations_year' THEN ((crossref->>'is-referenced-by-count')::double precision / (date_part('year', now()) - (paper.crossref->'published'->'date-parts'->0->>0)::double precision + 0.1))::double precision
    WHEN 'citations' THEN (crossref->>'is-referenced-by-count')::double precision
  END * CASE WHEN $5 = 'desc' THEN -1 ELSE 1 END
LIMIT $1
`

type SearchPaperByTitleParams struct {
	Limit     int32       `json:"limit"`
	Title     string      `json:"title"`
	Author    string      `json:"author"`
	OrderBy   string      `json:"order_by"`
	Direction interface{} `json:"direction"`
}

type SearchPaperByTitleRow struct {
	ID            int64       `json:"id"`
	Title         string      `json:"title"`
	Doi           string      `json:"doi"`
	Url           pgtype.Text `json:"url"`
	Journal       string      `json:"journal"`
	Author        interface{} `json:"author"`
	Published     pgtype.Date `json:"published"`
	Citations     interface{} `json:"citations"`
	CitationsYear float64     `json:"citations_year"`
}

func (q *Queries) SearchPaperByTitle(ctx context.Context, arg SearchPaperByTitleParams) ([]SearchPaperByTitleRow, error) {
	rows, err := q.db.Query(ctx, searchPaperByTitle,
		arg.Limit,
		arg.Title,
		arg.Author,
		arg.OrderBy,
		arg.Direction,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchPaperByTitleRow
	for rows.Next() {
		var i SearchPaperByTitleRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Doi,
			&i.Url,
			&i.Journal,
			&i.Author,
			&i.Published,
			&i.Citations,
			&i.CitationsYear,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
