-- +migrate Up
create table activities (
    id              integer         not null    auto_increment,
    day             date            not null,
    user_id         integer         not null,
    project_id      integer         not null,
    duration        integer         not null,
    primary key (id),
    index (day),
    index (user_id),
    index (project_id),
    foreign key (user_id) references users(id) on delete restrict,
    foreign key (project_id) references projects(id) on delete restrict
);

-- +migrate Down
drop table activities;
