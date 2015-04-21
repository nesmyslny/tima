-- +migrate Up
create table project (
    id                  integer         not null    auto_increment,
    project_category_id integer         not null,
    ref_id              varchar(5)      not null,
    ref_id_complete     varchar(30)     not null,
    responsible_user_id integer         not null,
    manager_user_id     integer         not null,
    title               varchar(100)    not null,
    description         text                null,
    version             integer         not null,
    primary key (id),
    index idx_project_project_category_id (project_category_id),
    index idx_project_ref_id (ref_id),
    unique index idx_project_ref_id_complete (ref_id_complete),
    index idx_project_responsible_user_id (responsible_user_id),
    index idx_project_manager_user_id (manager_user_id),
    constraint fk_project_project_category foreign key (project_category_id) references project_category(id) on delete restrict,
    constraint fk_project_user_responsible foreign key (responsible_user_id) references user(id) on delete restrict,
    constraint fk_project_user_manager foreign key (manager_user_id) references user(id) on delete restrict
);

-- +migrate Down
drop table project;
