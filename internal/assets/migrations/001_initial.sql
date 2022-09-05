-- +migrate Up
create table blobs (
    id bigserial primary key,
    information jsonb not null
);

-- +migrate Down
drop table blobs;
