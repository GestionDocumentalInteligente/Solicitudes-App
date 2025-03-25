-- +goose Up
-- Table to store account information
CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    account_number INTEGER NOT NULL UNIQUE
);

-- Table to store real estate information
CREATE TABLE IF NOT EXISTS real_estate (
    real_estate_id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(account_id),
    street VARCHAR(255),
    number INTEGER,
    locality VARCHAR(255),
    floor INTEGER,
    apartment VARCHAR(10)
);

-- Table to store parcel details
CREATE TABLE IF NOT EXISTS parcels (
    parcel_id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(account_id),
    parcel_number INTEGER,
    parcel_letter VARCHAR(10),
    polygon INTEGER,
    functional_unit INTEGER,
    complementary_unit INTEGER
);

-- Table to store owner information
CREATE TABLE IF NOT EXISTS owners (
    owner_id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(account_id),
    primary_owner VARCHAR(255),
    secondary_owner VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS owners;
DROP TABLE IF EXISTS parcels;
DROP TABLE IF EXISTS real_estate;
DROP TABLE IF EXISTS accounts;
