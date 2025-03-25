-- +goose Up
-- Table to store account information
CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    account_number BIGINT NOT NULL UNIQUE  -- Se utiliza BIGINT para permitir un rango amplio de números de cuenta
);

-- Table to store real estate information
CREATE TABLE IF NOT EXISTS real_estate (
    real_estate_id SERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(account_id),  -- Relación con la tabla accounts
    street VARCHAR(255),                                -- Tamaño adecuado para direcciones
    number SMALLINT CHECK (number > 0),                 -- Número de calle, utilizando SMALLINT ya que típicamente no es un valor muy alto
    locality VARCHAR(100),                              -- Tamaño reducido, suficiente para localidades
    floor SMALLINT,                                     -- Se utiliza SMALLINT para almacenar números de piso
    apartment VARCHAR(10)                               -- Suficiente para el formato común de apartamentos, como "A" o "12B"
);

-- Table to store parcel details
CREATE TABLE IF NOT EXISTS parcels (
    parcel_id SERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(account_id),  -- Relación con la tabla accounts
    parcel_number BIGINT CHECK (parcel_number > 0),     -- Utilizando BIGINT para permitir un rango amplio de números de parcela
    parcel_letter CHAR(5),                              -- Se utiliza CHAR(5) para una única letra o combinaciones pequeñas
    polygon SMALLINT,                                   -- Número de polígono; SMALLINT es suficiente en la mayoría de los casos
    functional_unit SMALLINT,                           -- Unidad funcional, usando SMALLINT para limitar el tamaño
    complementary_unit SMALLINT                         -- Unidad complementaria, usando SMALLINT para limitar el tamaño
);

-- Table to store owner information
CREATE TABLE IF NOT EXISTS owners (
    owner_id SERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(account_id),  -- Relación con la tabla accounts
    primary_owner VARCHAR(150),                         -- Ajustado para nombres completos largos
    secondary_owner VARCHAR(150)                        -- Ajustado para nombres completos largos
);

-- +goose Down
DROP TABLE IF EXISTS owners;
DROP TABLE IF EXISTS parcels;
DROP TABLE IF EXISTS real_estate;
DROP TABLE IF EXISTS accounts;
