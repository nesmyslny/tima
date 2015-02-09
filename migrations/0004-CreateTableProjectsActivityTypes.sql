-- +migrate Up
create table projects_activity_types (
    project_id          integer not null,
    activity_type_id    integer not null,
    primary key (project_id, activity_type_id),
    index(project_id),
    index(activity_type_id),
    constraint fk_projects_activity_types_projects foreign key (project_id) references projects(id) on delete cascade,
    constraint fk_projects_activity_types_activity_types foreign key (activity_type_id) references activity_types(id) on delete cascade
);

-- +migrate Down
drop table projects_activity_types;
