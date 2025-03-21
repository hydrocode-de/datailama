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

CREATE TABLE IF NOT EXISTS ollama_vector_collections
(
    id bigint primary key,
    name text not null
);

CREATE TABLE IF NOT EXISTS ollama_vectors
(
    id bigint primary key,
    collection_id bigint not null,
    text text not null,
    embedding vector(768) not null
);

