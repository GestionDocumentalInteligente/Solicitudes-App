-- +goose Up
-- Crear tabla activities
CREATE TABLE IF NOT EXISTS activities (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY (INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1),
    name character varying(255) NOT NULL,
    CONSTRAINT activities_pkey PRIMARY KEY (id)
);

-- Crear tabla intermedia request_activities
CREATE TABLE IF NOT EXISTS request_activities (
    request_id bigint NOT NULL,
    activity_id bigint NOT NULL,
    CONSTRAINT request_activities_pkey PRIMARY KEY (request_id, activity_id),
    CONSTRAINT request_activities_request_id_fkey FOREIGN KEY (request_id)
        REFERENCES requests(id)
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT request_activities_activity_id_fkey FOREIGN KEY (activity_id)
        REFERENCES activities(id)
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- Crear índices para mejorar el rendimiento
CREATE INDEX idx_request_activities_request_id ON request_activities(request_id);
CREATE INDEX idx_request_activities_activity_id ON request_activities(activity_id);

-- Eliminar la columna selected_activity de requests
ALTER TABLE requests DROP COLUMN IF EXISTS selected_activity;

-- +goose Down
-- Revertir los cambios
ALTER TABLE requests ADD COLUMN selected_activity bigint;

DROP TABLE IF EXISTS request_activities;
DROP TABLE IF EXISTS activities;