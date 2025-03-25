-- +goose Up
ALTER TABLE public.documents
ADD COLUMN name character varying(255) COLLATE pg_catalog."default";

-- +goose Down
ALTER TABLE public.documents
DROP COLUMN name;