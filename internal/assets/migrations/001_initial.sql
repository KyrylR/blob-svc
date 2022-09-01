-- +migrate Up

create table blobs (
    id bigserial primary key,
    data text not null
);

-- +migrate Down

drop table blobs