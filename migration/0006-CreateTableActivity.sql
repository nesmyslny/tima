-- +migrate Up
create table activity (
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
    constraint fk_activity_user foreign key (user_id) references user(id) on delete restrict,
    constraint fk_activity_project foreign key (project_id) references project(id) on delete restrict,
    constraint fk_activity_activity_type foreign key (activity_type_id) references activity_type(id) on delete restrict
);

-- +migrate Down
drop table activity;
