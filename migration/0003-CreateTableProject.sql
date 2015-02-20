-- +migrate Up
create table project (
    id                  integer         not null    auto_increment,
    project_category_id integer         not null,
    ref_id              varchar(5)      not null,
    ref_id_complete     varchar(30)     not null,
    title               varchar(100)    not null,
    primary key (id),
    index idx_project_project_category_id (project_category_id),
    index idx_project_ref_id (ref_id),
    unique index idx_project_ref_id_complete (ref_id_complete),
    constraint fk_project_project_category foreign key (project_category_id) references project_category(id) on delete restrict
);

-- +migrate Down
drop table project;
