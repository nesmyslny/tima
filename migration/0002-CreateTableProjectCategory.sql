-- +migrate Up
create table project_category (
    id              integer         not null    auto_increment,
    parent_id       integer         null,
    ref_id          varchar(5)      not null,
    ref_id_complete varchar(25)     not null,
    title           varchar(100)    not null,
    primary key (id),
    index idx_project_category_parent_id (parent_id),
    index idx_project_category_ref_id (ref_id),
    unique index idx_project_category_ref_id_complete (ref_id_complete),
    constraint fk_project_category_project_category foreign key (parent_id) references project_category(id) on delete cascade
);

-- +migrate Down
drop table project_category;
