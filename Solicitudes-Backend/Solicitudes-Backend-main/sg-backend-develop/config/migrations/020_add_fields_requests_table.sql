-- +goose Up
ALTER TABLE public.requests
ADD COLUMN abl_debt varchar(255),
ADD COLUMN selected_activity bigint,
ADD COLUMN estimated_time bigint,
ADD COLUMN insurance boolean DEFAULT false;

-- +goose Down
ALTER TABLE public.requests
DROP COLUMN abl_debt,
DROP COLUMN selected_activity,
DROP COLUMN estimated_time,
DROP COLUMN insurance;