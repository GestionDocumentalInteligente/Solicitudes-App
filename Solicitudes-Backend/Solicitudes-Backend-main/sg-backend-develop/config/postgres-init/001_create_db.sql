\connect template1;

-- Create the database 'sg_si_db' if it does not already exist.
DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sg_si_db') THEN
        CREATE DATABASE sg_si_db;
    END IF;
END
$$;

