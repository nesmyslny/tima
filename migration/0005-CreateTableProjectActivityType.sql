-- +migrate Up
create table project_activity_type (
    project_id          integer not null,
    activity_type_id    integer not null,
    primary key (project_id, activity_type_id),
    index(project_id),
    index(activity_type_id),
    constraint fk_project_activity_type_project foreign key (project_id) references project(id) on delete cascade,
    constraint fk_project_activity_type_activity_type foreign key (activity_type_id) references activity_type(id) on delete cascade
);

-- +migrate Down
drop table project_activity_type;
