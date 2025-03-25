-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Verificamos si existe la columna address_id
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'requests' 
        AND column_name = 'address_id'
    ) THEN
        -- Eliminamos la foreign key si existe
        ALTER TABLE requests 
        DROP CONSTRAINT IF EXISTS requests_address_id_fkey;

        -- Renombramos la columna
        ALTER TABLE requests 
        RENAME COLUMN address_id TO property_id;

        -- Agregamos la nueva foreign key si no existe
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'requests_property_id_fkey'
        ) THEN
            ALTER TABLE requests
            ADD CONSTRAINT requests_property_id_fkey 
            FOREIGN KEY (property_id)
            REFERENCES properties(property_id)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION;
        END IF;
    ELSIF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'requests' 
        AND column_name = 'property_id'
    ) THEN
        -- La columna ya está renombrada, solo verificamos la foreign key
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'requests_property_id_fkey'
        ) THEN
            ALTER TABLE requests
            ADD CONSTRAINT requests_property_id_fkey 
            FOREIGN KEY (property_id)
            REFERENCES properties(property_id)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION;
        END IF;
    ELSE
        RAISE NOTICE 'Ni address_id ni property_id existen en la tabla requests';
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    -- Verificamos si existe la columna property_id
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'requests' 
        AND column_name = 'property_id'
    ) THEN
        -- Eliminamos la foreign key si existe
        ALTER TABLE requests 
        DROP CONSTRAINT IF EXISTS requests_property_id_fkey;

        -- Renombramos la columna de vuelta
        ALTER TABLE requests 
        RENAME COLUMN property_id TO address_id;

        -- Verificamos si la tabla addresses existe antes de crear la foreign key
        IF EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'addresses'
        ) AND NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'requests_address_id_fkey'
        ) THEN
            ALTER TABLE requests
            ADD CONSTRAINT requests_address_id_fkey 
            FOREIGN KEY (address_id)
            REFERENCES addresses(id)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION;
        END IF;
    END IF;
END
$$;
-- +goose StatementEnd