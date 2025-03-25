-- 20240319000001_add_unaccent_extension.sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS unaccent SCHEMA public;

-- +goose Down
DROP EXTENSION IF EXISTS unaccent;