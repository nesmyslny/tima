-- +migrate Up
create table project (
    id                  integer         not null    auto_increment,
    project_category_id integer         not null,
    title               varchar(100)    not null,
    primary key (id),
    index idx_project_project_category_id (project_category_id),
    constraint fk_project_project_category foreign key (project_category_id) references project_category(id) on delete restrict
);

-- +migrate Down
drop table project;
