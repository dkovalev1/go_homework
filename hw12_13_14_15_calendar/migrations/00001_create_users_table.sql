-- +goose Up
CREATE TABLE event (
	id          text PRIMARY KEY,
	title       text,
	starttime   time,
	duration    int, -- in microseconds
	description text,
	userid      int,
	NotifyTime  int -- in microseconds
);

-- +goose Down
drop table events;
