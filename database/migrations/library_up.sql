-- ALTER TABLE books ADD files varchar(255) DEFAULT NULL;
create table books (
    id serial primary KEY,
    category_id bigint NULL REFERENCES categories,
    author_id bigint NULL REFERENCES authors,
    name varchar(255) NOT NULL,
    page varchar(255) DEFAULT NULL,
    files varchar(255) DEFAULT NULL
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