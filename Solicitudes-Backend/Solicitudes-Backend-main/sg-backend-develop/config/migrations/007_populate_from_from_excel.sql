-- +goose Up

-- Insert data into `accounts` table from `from_excel.inmuebles`
INSERT INTO accounts (account_number)
SELECT DISTINCT cuenta
FROM dblink('dbname=from_excel user=admin password=admin', 'SELECT cuenta FROM inmuebles')
AS from_excel_inmuebles(cuenta BIGINT);

-- Insert data into `real_estate` table from `from_excel.inmuebles`
INSERT INTO real_estate (account_id, street, number, locality, floor, apartment)
SELECT a.account_id,
       i.calle_inmueble,
       CASE WHEN i.altura_inmueble ~ '^[0-9]+$' THEN i.altura_inmueble::SMALLINT ELSE NULL END,  -- Verifica si es numérico
       i.localidad_inmueble,
       CASE WHEN i.piso_inmueble ~ '^[0-9]+$' THEN i.piso_inmueble::SMALLINT ELSE NULL END,      -- Verifica si es numérico
       i.depto_inmueble
FROM dblink('dbname=from_excel user=admin password=your_password', 
            'SELECT cuenta, calle_inmueble, altura_inmueble, localidad_inmueble, piso_inmueble, depto_inmueble FROM inmuebles') 
AS i(cuenta BIGINT, calle_inmueble VARCHAR(100), altura_inmueble VARCHAR, localidad_inmueble VARCHAR(100), piso_inmueble VARCHAR, depto_inmueble VARCHAR(10))
JOIN accounts a ON i.cuenta = a.account_number;

-- Insert data into `parcels` table from `from_excel.inmuebles`
INSERT INTO parcels (account_id, parcel_number, parcel_letter, polygon, functional_unit, complementary_unit)
SELECT a.account_id,
       CASE WHEN i.parcela ~ '^[0-9]+$' THEN i.parcela::BIGINT ELSE NULL END,
       i.letraparcela,
       CASE WHEN i.poligono ~ '^[0-9]+$' THEN i.poligono::SMALLINT ELSE NULL END,
       CASE WHEN i.unidfuncional ~ '^[0-9]+$' THEN i.unidfuncional::SMALLINT ELSE NULL END,
       CASE WHEN i.unidcomplem ~ '^[0-9]+$' THEN i.unidcomplem::SMALLINT ELSE NULL END
FROM dblink('dbname=from_excel user=admin password=your_password', 
            'SELECT cuenta, parcela, letraparcela, poligono, unidfuncional, unidcomplem FROM inmuebles') 
AS i(cuenta BIGINT, parcela VARCHAR, letraparcela CHAR(5), poligono VARCHAR, unidfuncional VARCHAR, unidcomplem VARCHAR)
JOIN accounts a ON i.cuenta = a.account_number;

-- Insert data into `owners` table from `from_excel.inmuebles`
INSERT INTO owners (account_id, primary_owner, secondary_owner)
SELECT a.account_id, i.titular, i.titularadjunto
FROM dblink('dbname=from_excel user=admin password=your_password', 
            'SELECT cuenta, titular, titularadjunto FROM inmuebles') 
AS i(cuenta BIGINT, titular VARCHAR(150), titularadjunto VARCHAR(150))
JOIN accounts a ON i.cuenta = a.account_number;

-- +goose Down
DELETE FROM owners;
DELETE FROM parcels;
DELETE FROM real_estate;
DELETE FROM accounts;
