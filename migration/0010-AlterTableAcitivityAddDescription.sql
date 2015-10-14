-- +migrate Up
alter table activity add column description text null after duration;

-- +migrate Down
alter table activity drop column description;
