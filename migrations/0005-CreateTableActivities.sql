-- +migrate Up
create table activities (
    id                  integer         not null    auto_increment,
    day                 date            not null,
    user_id             integer         not null,
    project_id          integer         not null,
    activity_type_id    integer         not null,
    duration            integer         not null,
    primary key (id),
    index (day),
    index (user_id),
    index (project_id),
    constraint fk_activities_users foreign key (user_id) references users(id) on delete restrict,
    constraint fk_activities_projects foreign key (project_id) references projects(id) on delete restrict,
    constraint fk_activities_activity_types foreign key (activity_type_id) references activity_types(id) on delete restrict
);

-- +migrate Down
drop table activities;
