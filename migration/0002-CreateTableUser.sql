-- +migrate Up
create table user (
    id              integer         not null    auto_increment,
    role            integer         not null,
    department_id   integer         not null,
    username        varchar(50)     not null,
    password_hash   binary(60)      not null,
    first_name      varchar(50)         null,
    last_name       varchar(50)         null,
    email           varchar(254)    not null,
    version         integer         not null,
    primary key (id),
    unique index idx_user_username (username)
);

-- +migrate Down
drop table user;
