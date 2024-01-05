create table books (
    id serial primary KEY,
    name varchar(255) NOT NULL,
    page bigint DEFAULT NULL,
    category_id bigint DEFAULT NULL REFERENCES categories ON DELETE CASCADE,
    author_id bigint DEFAULT NULL REFERENCES authors ON DELETE CASCADE
);

create table categories (
    id serial primary KEY,
    name varchar(255) NOT NULL
);

create table authors (
    id serial primary KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) DEFAULT NULL
);
