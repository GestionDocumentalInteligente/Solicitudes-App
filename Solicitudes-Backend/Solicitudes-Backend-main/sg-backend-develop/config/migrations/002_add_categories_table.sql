-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Check if the seed category already exists before inserting
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Planeamiento Urbano y Obras Particulares') THEN
        INSERT INTO categories (name) 
        VALUES ('Planeamiento Urbano y Obras Particulares');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM categories WHERE name = 'Transito y Movilidad') THEN
        INSERT INTO categories (name, description, is_active) 
        VALUES ('Transito y Movilidad', 'Servicios de transporte público, estacionamiento medido y seguridad vial.', false);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
