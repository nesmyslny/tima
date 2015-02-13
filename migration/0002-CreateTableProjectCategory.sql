-- +migrate Up
create table project_category (
    id              integer         not null    auto_increment,
    parent_id       integer         null,
    title           varchar(100)    not null,
    primary key (id)
);

-- +migrate Down
drop table project_category;
