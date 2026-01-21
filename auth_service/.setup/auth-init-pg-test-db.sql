-- Create test user
CREATE USER auth_test_db_user WITH PASSWORD 'changeit';

-- Create test database
CREATE DATABASE auth_test_db OWNER auth_test_db_user;

GRANT ALL PRIVILEGES ON DATABASE auth_test_db TO auth_test_db_user;