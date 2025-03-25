-- +goose Up
-- +goose StatementBegin

-- Primero, eliminamos todas las foreign keys para evitar problemas de dependencias
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT constraint_name, table_name 
              FROM information_schema.table_constraints 
              WHERE constraint_type = 'FOREIGN KEY' 
              AND table_schema = current_schema()
              AND table_name IN ('users', 'user_roles', 'role_permissions'))
    LOOP
        EXECUTE 'ALTER TABLE ' || quote_ident(r.table_name) || ' DROP CONSTRAINT ' || quote_ident(r.constraint_name);
    END LOOP;
END $$;

-- Ahora eliminamos los índices
DROP INDEX IF EXISTS idx_persons_cuil;
DROP INDEX IF EXISTS idx_persons_email;

-- Eliminamos constraints de unicidad
DO $$ 
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT constraint_name, table_name 
              FROM information_schema.table_constraints 
              WHERE constraint_type = 'UNIQUE' 
              AND table_schema = current_schema()
              AND table_name IN ('persons'))
    LOOP
        EXECUTE 'ALTER TABLE ' || quote_ident(r.table_name) || ' DROP CONSTRAINT IF EXISTS ' || quote_ident(r.constraint_name);
    END LOOP;
END $$;

-- Ahora podemos modificar las tablas de manera segura

-- Modificar persons
ALTER TABLE persons DROP CONSTRAINT IF EXISTS persons_pkey CASCADE;
ALTER TABLE persons 
    DROP COLUMN IF EXISTS uuid CASCADE,
    ALTER COLUMN cuil TYPE VARCHAR(20),
    ALTER COLUMN cuil SET NOT NULL,
    ALTER COLUMN dni TYPE VARCHAR(20),
    ALTER COLUMN dni SET NOT NULL,
    ALTER COLUMN first_name TYPE VARCHAR(255),
    ALTER COLUMN first_name SET NOT NULL,
    ALTER COLUMN last_name TYPE VARCHAR(255),
    ALTER COLUMN last_name SET NOT NULL,
    ALTER COLUMN email TYPE VARCHAR(255),
    ALTER COLUMN phone TYPE VARCHAR(50),
    ADD COLUMN IF NOT EXISTS id bigint NOT NULL GENERATED ALWAYS AS IDENTITY;

-- Modificar users
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;
ALTER TABLE users 
    DROP COLUMN IF EXISTS uuid CASCADE,
    ADD COLUMN IF NOT EXISTS id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    ADD COLUMN IF NOT EXISTS person_id bigint,
    ADD COLUMN IF NOT EXISTS email_validated BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS accepts_notifications BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS logged_at TIMESTAMPTZ DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;

-- Modificar user_roles
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS user_roles_pkey CASCADE;
ALTER TABLE user_roles 
    DROP COLUMN IF EXISTS user_uuid CASCADE,
    ADD COLUMN IF NOT EXISTS user_id bigint,
    ADD COLUMN IF NOT EXISTS role_id INT;

-- Eliminar columnas extras de cada tabla
DO $$ 
DECLARE
    col_name text;
BEGIN
    -- Limpiar persons
    FOR col_name IN (
        SELECT column_name 
        FROM information_schema.columns 
        WHERE table_name = 'persons' 
        AND column_name NOT IN ('id', 'cuil', 'dni', 'first_name', 'last_name', 'email', 'phone')
    )
    LOOP
        EXECUTE 'ALTER TABLE persons DROP COLUMN IF EXISTS ' || quote_ident(col_name) || ' CASCADE';
    END LOOP;

    -- Limpiar users
    FOR col_name IN (
        SELECT column_name 
        FROM information_schema.columns 
        WHERE table_name = 'users' 
        AND column_name NOT IN ('id', 'person_id', 'email_validated', 'accepts_notifications', 
                              'created_at', 'logged_at', 'updated_at', 'deleted_at')
    )
    LOOP
        EXECUTE 'ALTER TABLE users DROP COLUMN IF EXISTS ' || quote_ident(col_name) || ' CASCADE';
    END LOOP;

    -- Limpiar user_roles
    FOR col_name IN (
        SELECT column_name 
        FROM information_schema.columns 
        WHERE table_name = 'user_roles' 
        AND column_name NOT IN ('user_id', 'role_id')
    )
    LOOP
        EXECUTE 'ALTER TABLE user_roles DROP COLUMN IF EXISTS ' || quote_ident(col_name) || ' CASCADE';
    END LOOP;
END $$;

-- Reestablecer primary keys
ALTER TABLE persons ADD PRIMARY KEY (id);
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE user_roles ADD PRIMARY KEY (user_id, role_id);

-- Reestablecer unique constraints
ALTER TABLE persons ADD CONSTRAINT persons_cuil_key UNIQUE (cuil);
ALTER TABLE persons ADD CONSTRAINT persons_email_key UNIQUE (email);

-- Reestablecer foreign keys
ALTER TABLE users 
    ADD CONSTRAINT fk_person 
    FOREIGN KEY (person_id) 
    REFERENCES persons(id) 
    ON DELETE SET NULL;

ALTER TABLE user_roles 
    ADD CONSTRAINT fk_user_roles_user 
    FOREIGN KEY (user_id) 
    REFERENCES users(id) 
    ON DELETE CASCADE,
    ADD CONSTRAINT fk_user_roles_role 
    FOREIGN KEY (role_id) 
    REFERENCES roles(id) 
    ON DELETE CASCADE;

ALTER TABLE role_permissions 
    ADD CONSTRAINT fk_role_permissions_role 
    FOREIGN KEY (role_id) 
    REFERENCES roles(id) 
    ON DELETE CASCADE,
    ADD CONSTRAINT fk_role_permissions_permission 
    FOREIGN KEY (permission_id) 
    REFERENCES permissions(id) 
    ON DELETE CASCADE;

-- Recrear índices
CREATE INDEX idx_persons_cuil ON persons(cuil);
CREATE INDEX idx_persons_email ON persons(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- No implementamos down para evitar pérdida de datos
-- +goose StatementEnd