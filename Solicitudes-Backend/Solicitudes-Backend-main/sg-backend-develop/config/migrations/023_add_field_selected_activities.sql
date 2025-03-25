-- +goose Up
-- Agregar la columna selected_activities como un array de integers
ALTER TABLE requests ADD COLUMN selected_activities integer[] DEFAULT '{}';

-- +goose Down
-- Eliminar la columna en caso de rollback
ALTER TABLE requests DROP COLUMN IF EXISTS selected_activities;