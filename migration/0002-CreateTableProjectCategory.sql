-- +migrate Up
create table project_category (
    id              integer         not null    auto_increment,
    parent_id       integer         null,
    title           varchar(100)    not null,
    primary key (id),
    index idx_project_category_parent_id (parent_id),
    constraint fk_project_category_project_category foreign key (parent_id) references project_category(id) on delete cascade
);

-- +migrate Down
drop table project_category;
