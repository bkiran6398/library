-- +goose Up
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    isbn TEXT NOT NULL UNIQUE,
    published_year INT,
    copies_total INT NOT NULL CHECK (copies_total >= 0),
    copies_available INT NOT NULL CHECK (copies_available >= 0 AND copies_available <= copies_total),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_books_title ON books (title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books (author);
CREATE INDEX IF NOT EXISTS idx_books_created_at ON books (created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS books;



