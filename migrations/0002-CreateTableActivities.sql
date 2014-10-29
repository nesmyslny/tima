-- +migrate Up
create table activities (
    id              integer         not null    auto_increment,
    day             date            not null,
    user_id         integer         not null,
    text            varchar(100)    not null,
    duration        integer         not null,
    primary key (id),
    index (day),
    index (user_id),
    foreign key (user_id) references users(id) on delete restrict
);

-- +migrate Down
drop table activities;
