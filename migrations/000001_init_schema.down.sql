-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_service_providers_updated_at ON service_providers;
DROP TRIGGER IF EXISTS update_services_updated_at ON services;
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS provider_services;
DROP TABLE IF EXISTS service_providers;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS user_role;

