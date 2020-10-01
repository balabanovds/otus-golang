-- +goose Up
-- SQL in this section is executed when the migration is applied.
create extension if not exists "uuid-ossp";

create table events (
    id serial primary key,
    title varchar(256) not null,
    start_at timestamptz,
    end_at timestamptz,
    description text,
    user_id int,
    remind_at timestamptz,
    created_at timestamp default now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists events;
