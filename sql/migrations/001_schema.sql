-- +goose Up
CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS themes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);
CREATE TABLE IF NOT EXISTS book_ideas (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    submitter TEXT NOT NULL,
    theme_id INTEGER REFERENCES themes(id),
    genre_id INTEGER REFERENCES genres(id)
);

-- +goose Down
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS book_ideas;
