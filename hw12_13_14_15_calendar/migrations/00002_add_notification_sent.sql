-- +goose Up

ALTER TABLE event 
    ADD COLUMN notification_sent boolean DEFAULT false,
    ALTER COLUMN starttime TYPE timestamp with time zone USING current_date + starttime;

-- +goose Down

ALTER TABLE event 
    ALTER COLUMN starttime TYPE time,
    DROP COLUMN notification_sent;
