-- Creates the petstore application role.
-- The migration runner must SET migration.petstore_password
-- to the value of PETSTORE_PASSWORD before running this.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE rolname = 'petstore'
    ) THEN
        EXECUTE format(
            'CREATE ROLE petstore LOGIN PASSWORD %L',
            current_setting('migration.petstore_password')
        );
    END IF;
END
$$;
