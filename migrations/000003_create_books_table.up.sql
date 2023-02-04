CREATE TABLE IF NOT EXISTS books (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    
    title text not null,
    author text not null,
    year integer not null,
    language text not null,
    price integer not null,
    quantity integer not null,
    -- genres column is array of zero-or-more text values. 
    genres text[] not NULL,
    version integer NOT NULL DEFAULT 1
);