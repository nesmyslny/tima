-- +migrate Up
create table project_department (
    project_id      integer not null,
    department_id   integer not null,
    primary key (project_id, department_id),
    index idx_project_department_project_id (project_id),
    index idx_project_department_department_id (department_id),
    constraint fk_project_department_project foreign key (project_id) references project(id) on delete cascade,
    constraint fk_project_department_department foreign key (department_id) references department(id) on delete cascade
);

-- +migrate Down
drop table project_department;
