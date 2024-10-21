-- +goose Up
create table if not exists users (
    id serial primary key,
    username text not null unique CHECK (length(username) <= 20),
    password text not null
);

create table if not exists assets (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    name text CHECK (length(name) > 0),
    descr text,
    price int not null
);

-- +goose Down
drop table if exists users;
drop table if exists assets;