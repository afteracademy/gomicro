-- Create test user
CREATE USER goserver_test_db_user WITH PASSWORD 'changeit';

-- Create test database
CREATE DATABASE goserver_test_db OWNER goserver_test_db_user;

GRANT ALL PRIVILEGES ON DATABASE goserver_test_db TO goserver_test_db_user;