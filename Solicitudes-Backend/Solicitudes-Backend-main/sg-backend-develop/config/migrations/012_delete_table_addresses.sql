-- +goose Up
-- +goose StatementBegin
-- Primero eliminamos cualquier foreign key que apunte a la tabla addresses
DO $$
BEGIN
    -- Buscamos y eliminamos todas las foreign keys que referencian a addresses
    EXECUTE (
        SELECT string_agg('ALTER TABLE ' || quote_ident(tc.table_schema) || '.' || quote_ident(tc.table_name) || 
               ' DROP CONSTRAINT ' || quote_ident(tc.constraint_name), '; ')
        FROM information_schema.table_constraints tc
        JOIN information_schema.constraint_column_usage ccu 
        ON ccu.constraint_name = tc.constraint_name
        WHERE tc.constraint_type = 'FOREIGN KEY' 
        AND ccu.table_name = 'addresses'
    );
EXCEPTION 
    WHEN OTHERS THEN NULL; -- Ignoramos errores si no hay foreign keys
END $$;

-- Eliminamos triggers si existen
DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;

-- Eliminamos índices si existen
DROP INDEX IF EXISTS idx_addresses_street;
DROP INDEX IF EXISTS idx_addresses_street_number;
DROP INDEX IF EXISTS idx_addresses_postal_code;
DROP INDEX IF EXISTS idx_addresses_city;
DROP INDEX IF EXISTS idx_addresses_state;

-- Finalmente eliminamos la tabla
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Recreamos la tabla en caso de rollback
CREATE TABLE IF NOT EXISTS addresses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(20) NOT NULL,
    floor VARCHAR(20),
    apartment VARCHAR(20),
    functional_unit VARCHAR(50),
    postal_code VARCHAR(20),
    city VARCHAR(100),
    state VARCHAR(100),
    is_common_area BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recreamos el trigger de updated_at
CREATE OR REPLACE TRIGGER update_addresses_updated_at 
    BEFORE UPDATE ON addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd