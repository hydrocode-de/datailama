-- name: GetPaperStatistics :many
SELECT j.title,
    paper.issn,
    ((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0] AS year,
    count(*) AS count
   FROM paper
     JOIN journals j ON paper.issn = j.issn
  GROUP BY paper.issn, j.title, (((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0])
  ORDER BY (((paper.crossref -> 'published'::text) -> 'date-parts'::text)[0][0]);


-- name: SearchPaperByTitle :many
SELECT 
  paper.id,
  paper.title,
  paper.doi,
  paper.url,
  journals.title as journal,
  paper.crossref->'author'->0->>'family' || ', ' || (paper.crossref->'author'->0->>'given')::text as author,
  date(paper.crossref->'published'->'date-parts'->>0) as published,
  crossref->>'is-referenced-by-count' as citations,
  ((crossref->>'is-referenced-by-count')::numeric / (date_part('year', now()) - (paper.crossref->'published'->'date-parts'->0->>0)::numeric + 0.1))::numeric(10,2) as "citations_year"
FROM paper
JOIN journals ON journals.issn=paper.issn
WHERE paper.title ILIKE '%' || @title::text || '%'
AND paper.crossref->'author'->0->>'family' ILIKE '%' || @author::text || '%'
ORDER BY 
  CASE @order_by::text
    WHEN 'citations_year' THEN ((crossref->>'is-referenced-by-count')::numeric / (date_part('year', now()) - (paper.crossref->'published'->'date-parts'->0->>0)::numeric + 0.1))::numeric(10,2)
    WHEN 'citations' THEN (crossref->>'is-referenced-by-count')::numeric
  END * CASE WHEN @direction = 'desc' THEN -1 ELSE 1 END
LIMIT $1;