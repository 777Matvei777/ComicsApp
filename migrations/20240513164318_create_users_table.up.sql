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
    roles text not null
);


INSERT INTO users(username, pass, roles) VALUES ('Matvei', '1234', 'admin');
INSERT INTO users(username, pass, roles) VALUES ('Igor', '3456', 'user');
INSERT INTO users(username, pass, roles) VALUES ('Artem', '4567', 'user');

