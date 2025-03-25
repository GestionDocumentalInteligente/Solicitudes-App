-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Renombrar tabla real_estate a properties
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'real_estate') THEN
        ALTER TABLE real_estate RENAME TO properties;
        -- Solo intentamos renombrar la columna si la tabla existía y fue renombrada
        ALTER TABLE properties RENAME COLUMN real_estate_id TO property_id;
    END IF;

    -- Renombrar tabla accounts a abl
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'accounts') THEN
        ALTER TABLE accounts RENAME TO abl;
        -- Solo intentamos renombrar las columnas si la tabla existía y fue renombrada
        ALTER TABLE abl RENAME COLUMN account_id TO abl_id;
        ALTER TABLE abl RENAME COLUMN account_number TO abl_number;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    -- Revertir cambios de abl a accounts
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'abl') THEN
        ALTER TABLE abl RENAME COLUMN abl_id TO account_id;
        ALTER TABLE abl RENAME COLUMN abl_number TO account_number;
        ALTER TABLE abl RENAME TO accounts;
    END IF;

    -- Revertir cambios de properties a real_estate
    IF EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'properties') THEN
        ALTER TABLE properties RENAME COLUMN property_id TO real_estate_id;
        ALTER TABLE properties RENAME TO real_estate;
    END IF;
END
$$;
-- +goose StatementEnd