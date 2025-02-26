-- +goose Up
CREATE TABLE event (
	id          text PRIMARY KEY,
	title       text,
	starttime   time,
	duration    int, -- in seconds
	description text,
	userid      int,
	NotifyTime  int -- in seconds
);

-- +goose Down
drop table events;
