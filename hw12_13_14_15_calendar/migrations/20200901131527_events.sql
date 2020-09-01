-- +goose Up
-- SQL in this section is executed when the migration is applied.
create extension if not exists "uuid-ossp";

create table events (
    id uuid default uuid_generate_v4(),
    title varchar(256) not null,
    start_time timestamp,
    duration bigint,
    description text,
    user_id int,
    remind_duration bigint,
    created_at timestamp default now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists events;
