CREATE TABLE IF NOT EXISTS paper
(
    id bigint primary key,
    doi text NOT NULL,
    url text,
    issn text NOT NULL,
    title text NOT NULL,
    crossref jsonb NOT NULL,
    body text
);

CREATE TABLE IF NOT EXISTS journals 
(
    issn text primary key,
    title text not null,
    short text not null
);


