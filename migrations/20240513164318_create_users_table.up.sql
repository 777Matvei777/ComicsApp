CREATE TABLE comics(
    id SERIAL PRIMARY KEY,
    url text NOT NULL
);

CREATE TABLE keywords(
    id  serial primary key,
    keyword text not null,
    comic_id INTEGER REFERENCES comics(id)
);

CREATE TABLE keyword_index(
    keyword TEXT PRIMARY KEY,
    comic_ids INT[] NOT NULL
);

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    username text unique not null,
    pass text not null,
    role text not null
);
