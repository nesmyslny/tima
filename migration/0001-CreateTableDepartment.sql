-- +migrate Up
create table department (
    id              integer         not null    auto_increment,
    parent_id       integer         null,
    title           varchar(100)    not null,
    primary key (id),
    index idx_department_parent_id (parent_id),
    constraint fk_department_department foreign key (parent_id) references department(id) on delete cascade
);

-- +migrate Down
drop table department;
