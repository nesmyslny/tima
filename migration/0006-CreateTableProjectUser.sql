-- +migrate Up
create table project_user (
    project_id  integer not null,
    user_id     integer not null,
    primary key (project_id, user_id),
    index idx_project_user_project_id (project_id),
    index idx_project_user_user_id (user_id),
    constraint fk_project_user_project foreign key (project_id) references project(id) on delete cascade,
    constraint fk_project_user_user foreign key (user_id) references user(id) on delete cascade
);

-- +migrate Down
drop table project_user;
