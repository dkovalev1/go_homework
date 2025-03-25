-- +goose Up

CREATE TABLE notification ( --no primary key is intended
	eventid text,
	title   text,
	start   time,
	userid  int
);

-- +goose Down
drop table notification;
