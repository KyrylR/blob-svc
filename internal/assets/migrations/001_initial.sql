-- +migrate Up
create table blobs (
    id bigserial primary key,
    information text not null
);

-- +migrate Down
drop table blobs;
