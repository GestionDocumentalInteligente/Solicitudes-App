-- +goose Up
-- +goose StatementBegin

-- Función auxiliar para verificar si una tabla existe
CREATE OR REPLACE FUNCTION table_exists(t_name VARCHAR) RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_name = t_name
    );
END;
$$ LANGUAGE plpgsql;

-- Función auxiliar para verificar si una columna existe
CREATE OR REPLACE FUNCTION column_exists(t_name VARCHAR, c_name VARCHAR) RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT FROM information_schema.columns
        WHERE table_name = t_name AND column_name = c_name
    );
END;
$$ LANGUAGE plpgsql;

DO $$ 
BEGIN
    -- Verificar y crear/modificar tabla persons
    IF NOT table_exists('persons') THEN
        CREATE TABLE persons (
            id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
            cuil VARCHAR(20) NOT NULL UNIQUE,
            dni VARCHAR(20) NOT NULL,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE,
            phone VARCHAR(50),
            PRIMARY KEY (id)
        );
        
        -- Crear índices para nueva tabla
        CREATE INDEX idx_persons_cuil ON persons(cuil);
        CREATE INDEX idx_persons_email ON persons(email);
    ELSE
        -- Verificar y agregar columnas si no existen
        IF NOT column_exists('persons', 'cuil') THEN
            ALTER TABLE persons ADD COLUMN cuil VARCHAR(20) UNIQUE;
            CREATE INDEX idx_persons_cuil ON persons(cuil);
        END IF;
        
        IF NOT column_exists('persons', 'email') THEN
            ALTER TABLE persons ADD COLUMN email VARCHAR(255) UNIQUE;
            CREATE INDEX idx_persons_email ON persons(email);
        END IF;
        
        IF NOT column_exists('persons', 'phone') THEN
            ALTER TABLE persons ADD COLUMN phone VARCHAR(50);
        END IF;
    END IF;

    -- Verificar y crear tabla permissions
    IF NOT table_exists('permissions') THEN
        CREATE TABLE permissions (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) NOT NULL UNIQUE,
            description TEXT
        );
    END IF;

    -- Verificar y crear tabla roles
    IF NOT table_exists('roles') THEN
        CREATE TABLE roles (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) NOT NULL UNIQUE,
            description TEXT
        );
    END IF;

    -- Verificar y crear/modificar tabla users
    IF NOT table_exists('users') THEN
        CREATE TABLE users (
            id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
            person_id bigint,
            email_validated BOOLEAN DEFAULT FALSE,
            accepts_notifications BOOLEAN DEFAULT FALSE,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            logged_at TIMESTAMPTZ DEFAULT NULL,
            updated_at TIMESTAMPTZ DEFAULT NULL,
            deleted_at TIMESTAMPTZ DEFAULT NULL,
            PRIMARY KEY (id),
            CONSTRAINT fk_person FOREIGN KEY(person_id) 
                REFERENCES persons(id) ON DELETE SET NULL
        );
    ELSE
        -- Verificar y agregar columnas si no existen
        IF NOT column_exists('users', 'email_validated') THEN
            ALTER TABLE users ADD COLUMN email_validated BOOLEAN DEFAULT FALSE;
        END IF;
        
        IF NOT column_exists('users', 'accepts_notifications') THEN
            ALTER TABLE users ADD COLUMN accepts_notifications BOOLEAN DEFAULT FALSE;
        END IF;
        
        IF NOT column_exists('users', 'logged_at') THEN
            ALTER TABLE users ADD COLUMN logged_at TIMESTAMPTZ DEFAULT NULL;
        END IF;
    END IF;

    -- Verificar y crear tabla user_roles
    IF NOT table_exists('user_roles') THEN
        CREATE TABLE user_roles (
            user_id bigint,
            role_id INT,
            PRIMARY KEY(user_id, role_id),
            CONSTRAINT fk_user_roles_user 
                FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
            CONSTRAINT fk_user_roles_role 
                FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
        );
    END IF;

    -- Verificar y crear tabla role_permissions
    IF NOT table_exists('role_permissions') THEN
        CREATE TABLE role_permissions (
            role_id INT,
            permission_id INT,
            PRIMARY KEY(role_id, permission_id),
            CONSTRAINT fk_role_permissions_role 
                FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
            CONSTRAINT fk_role_permissions_permission 
                FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE
        );
    END IF;

END $$;

-- Limpiar funciones auxiliares
DROP FUNCTION IF EXISTS table_exists(VARCHAR);
DROP FUNCTION IF EXISTS column_exists(VARCHAR, VARCHAR);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- No se implementa el down ya que es una migración de actualización
-- y no queremos perder datos existentes

-- +goose StatementEnd