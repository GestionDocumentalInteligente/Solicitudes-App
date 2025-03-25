-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS request_types (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id INTEGER REFERENCES categories(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    requires_documentation BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Check if the seed request_type already exists before inserting
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM request_types WHERE name = 'Aviso de Obras') THEN
        INSERT INTO request_types (category_id, name) 
        VALUES (1, 'Aviso de Obras');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM request_types WHERE name = 'Defensa al consumidor') THEN
        INSERT INTO request_types (category_id, name, is_active) 
        VALUES (1, 'Defensa al consumidor', false);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM request_types WHERE name = 'Impacto ambiental') THEN
        INSERT INTO request_types (category_id, name, is_active) 
        VALUES (1, 'Impacto ambiental', false);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_types;
-- +goose StatementEnd
