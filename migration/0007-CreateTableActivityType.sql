-- +migrate Up
create table activity_type (
    id              integer         not null    auto_increment,
    title           varchar(100)    not null,
    version         integer         not null,
    primary key (id)
);

-- +migrate Down
drop table activity_type;
