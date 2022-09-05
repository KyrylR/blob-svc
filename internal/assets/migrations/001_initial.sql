-- +migrate Up
create table blobs (
    id bigserial primary key,
    information jsonb not null,
    owner_address character(64)
);

-- +migrate Down
drop table blobs;
