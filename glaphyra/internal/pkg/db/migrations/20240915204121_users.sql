-- +goose Up
create table users (
    id bigint not null primary key,
    first_name text not null,
    last_name text not null,
    middle_name text,
    created_at timestamp default now()
);
CREATE SEQUENCE users_id START 1;

-- +goose Down
drop table users;
drop sequence users_id;
