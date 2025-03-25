-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Properties: account_id -> abl_id
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'properties' 
        AND column_name = 'account_id'
    ) THEN
        ALTER TABLE properties RENAME COLUMN account_id TO abl_id;
    END IF;

    -- Parcels: account_id -> abl_id
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'parcels' 
        AND column_name = 'account_id'
    ) THEN
        ALTER TABLE parcels RENAME COLUMN account_id TO abl_id;
    END IF;

    -- Owners: account_id -> abl_id
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'owners' 
        AND column_name = 'account_id'
    ) THEN
        ALTER TABLE owners RENAME COLUMN account_id TO abl_id;
    END IF;

    -- Actualizar foreign keys
    -- Properties
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'properties' 
        AND column_name = 'abl_id'
    ) AND NOT EXISTS (
        SELECT FROM information_schema.table_constraints
        WHERE constraint_name = 'properties_abl_id_fkey'
    ) THEN
        ALTER TABLE properties
            DROP CONSTRAINT IF EXISTS properties_account_id_fkey,
            ADD CONSTRAINT properties_abl_id_fkey 
            FOREIGN KEY (abl_id) REFERENCES abl(abl_id);
    END IF;

    -- Parcels
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'parcels' 
        AND column_name = 'abl_id'
    ) AND NOT EXISTS (
        SELECT FROM information_schema.table_constraints
        WHERE constraint_name = 'parcels_abl_id_fkey'
    ) THEN
        ALTER TABLE parcels
            DROP CONSTRAINT IF EXISTS parcels_account_id_fkey,
            ADD CONSTRAINT parcels_abl_id_fkey 
            FOREIGN KEY (abl_id) REFERENCES abl(abl_id);
    END IF;

    -- Owners
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'owners' 
        AND column_name = 'abl_id'
    ) AND NOT EXISTS (
        SELECT FROM information_schema.table_constraints
        WHERE constraint_name = 'owners_abl_id_fkey'
    ) THEN
        ALTER TABLE owners
            DROP CONSTRAINT IF EXISTS owners_account_id_fkey,
            ADD CONSTRAINT owners_abl_id_fkey 
            FOREIGN KEY (abl_id) REFERENCES abl(abl_id);
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    -- Revertir foreign keys
    -- Properties
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'properties' 
        AND column_name = 'abl_id'
    ) THEN
        ALTER TABLE properties
            DROP CONSTRAINT IF EXISTS properties_abl_id_fkey;
        
        ALTER TABLE properties RENAME COLUMN abl_id TO account_id;
        
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'properties_account_id_fkey'
        ) THEN
            ALTER TABLE properties
            ADD CONSTRAINT properties_account_id_fkey 
            FOREIGN KEY (account_id) REFERENCES abl(abl_id);
        END IF;
    END IF;

    -- Parcels
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'parcels' 
        AND column_name = 'abl_id'
    ) THEN
        ALTER TABLE parcels
            DROP CONSTRAINT IF EXISTS parcels_abl_id_fkey;
        
        ALTER TABLE parcels RENAME COLUMN abl_id TO account_id;
        
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'parcels_account_id_fkey'
        ) THEN
            ALTER TABLE parcels
            ADD CONSTRAINT parcels_account_id_fkey 
            FOREIGN KEY (account_id) REFERENCES abl(abl_id);
        END IF;
    END IF;

    -- Owners
    IF EXISTS (
        SELECT FROM information_schema.columns 
        WHERE table_name = 'owners' 
        AND column_name = 'abl_id'
    ) THEN
        ALTER TABLE owners
            DROP CONSTRAINT IF EXISTS owners_abl_id_fkey;
        
        ALTER TABLE owners RENAME COLUMN abl_id TO account_id;
        
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints
            WHERE constraint_name = 'owners_account_id_fkey'
        ) THEN
            ALTER TABLE owners
            ADD CONSTRAINT owners_account_id_fkey 
            FOREIGN KEY (account_id) REFERENCES abl(abl_id);
        END IF;
    END IF;
END
$$;
-- +goose StatementEnd