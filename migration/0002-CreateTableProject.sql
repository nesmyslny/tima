-- +migrate Up
create table project (
    id              integer         not null    auto_increment,
    title           varchar(100)    not null,
    primary key (id)
);

-- +migrate Down
drop table project;
