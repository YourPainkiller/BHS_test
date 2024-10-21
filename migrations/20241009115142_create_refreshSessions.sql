-- +goose Up
create table if not exists refreshSessions (
    id serial primary key,
    user_id bigint,
    refresh_token text not null,
    fingerprint text not null,
    ip text not null,
    expires_in int,
    created_at text not null
);

-- +goose Down
drop table if exists refreshSessions;