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
  '' as match,
  0 as cosine_distance,
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
WHERE paper.title ILIKE '%' || @title::text || '%'
AND paper.crossref->'author'->0->>'family' ILIKE '%' || @author::text || '%'
ORDER BY 
  CASE @order_by::text
    WHEN 'citations_year' THEN ((crossref->>'is-referenced-by-count')::double precision / (date_part('year', now()) - (paper.crossref->'published'->'date-parts'->0->>0)::double precision + 0.1))::double precision
    WHEN 'citations' THEN (crossref->>'is-referenced-by-count')::double precision
  END * CASE WHEN @direction = 'desc' THEN -1 ELSE 1 END
LIMIT $1;

-- name: SearchPaperBody :many
WITH embed AS(
  SELECT @embedding::vector AS embedding
),
matches AS (
  SELECT 
    id,
    text,
    collection_id,
    embedding <=> (SELECT embed.embedding FROM embed) AS cosine_distance
  FROM ollama_vectors
  ORDER BY cosine_distance
  LIMIT $1
)
SELECT 
  m.text as match,
  m.cosine_distance,
  p.id,
  p.title,
  p.doi,
  p.url,
  j.title as journal,
  date(p.crossref->'published'->'date-parts'->>0) as published,
  p.crossref->>'is-referenced-by-count' as citations,
  ((p.crossref->>'is-referenced-by-count')::double precision / (date_part('year', now()) - (p.crossref->'published'->'date-parts'->0->>0)::double precision + 0.1))::double precision as "citations_year"
FROM matches m
JOIN ollama_vector_collections c ON m.collection_id=c.uuid
JOIN paper p ON p.doi=c.name
JOIN journals j ON p.issn=j.issn;