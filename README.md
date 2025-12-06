# Qasynda - Service Marketplace API

A REST API backend for a service marketplace platform that connects clients with service providers (plumbers, electricians, cleaners, etc.).

## Project Overview

Qasynda is a professional networking platform focused on local services. It allows:
- Users to register as Clients or Service Providers
- Service Providers to create profiles with skills, rates, and availability
- Clients to search and filter providers by service type, location, and rating
- Booking system for scheduling appointments
- Reviews & Ratings after service completion
- Real-time notifications for booking requests
- Chat service for direct communication between clients and providers

## Project Structure

```
qasynda/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/                  # Database models/entities
│   ├── repository/              # Database access layer
│   ├── service/                 # Business logic layer
│   ├── handler/                 # HTTP handlers (controllers)
│   ├── middleware/              # Auth, logging, rate limiting, CORS, recovery
│   ├── worker/                  # Background jobs (notifications)
│   └── config/                  # Configuration
├── migrations/                  # Database migrations (golang-migrate)
├── pkg/                         # Shared utilities
│   ├── auth/                    # JWT utilities
│   └── validator/               # Input validation
├── tests/                       # Test files
├── .env.example                 # Environment variables template
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

**Option 1: Local Development**
- Go 1.21 or higher
- PostgreSQL 12 or higher
- golang-migrate CLI tool (for database migrations)

**Option 2: Docker (Recommended)**
- Docker 20.10 or higher
- Docker Compose 2.0 or higher

## Setup Instructions

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd go-final
go mod download
```

### 2. Database Setup

Create a PostgreSQL database:

```bash
createdb qasynda
```

Or using psql:

```sql
CREATE DATABASE qasynda;
```

### 3. Environment Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
PORT=8080
ENV=development

DATABASE_URL=postgres://postgres:postgres@localhost:5432/qasynda?sslmode=disable
# OR use individual components:
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=qasynda
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

LOG_LEVEL=info
```

### 4. Run Database Migrations

Install golang-migrate:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# Or using Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/qasynda?sslmode=disable" up
```

To rollback:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/qasynda?sslmode=disable" down
```

### 5. Run the Application

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### 6. Build for Production

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

### 7. Docker Setup (Recommended)

#### Using Docker Compose (Production)

Build and run everything with Docker Compose:

```bash
# Build and start all services (API + PostgreSQL + Migrations)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

This will:
- Start PostgreSQL database
- Run database migrations automatically
- Build and start the API server

#### Using Docker Compose (Development with Hot Reload)

For development with automatic code reloading:

```bash
# Start services with hot reload
docker-compose -f docker-compose.dev.yml up

# Or run in background
docker-compose -f docker-compose.dev.yml up -d
```

The development setup includes:
- Hot reload using Air (automatically rebuilds on code changes)
- Volume mounting for live code updates
- Development-friendly environment variables

#### Manual Docker Build

```bash
# Build the Docker image
docker build -t qasynda-api .

# Run the container
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/qasynda?sslmode=disable" \
  -e JWT_SECRET="your-secret-key" \
  qasynda-api
```

#### Environment Variables for Docker

Create a `.env` file in the project root (optional, defaults are used):

```env
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h
LOG_LEVEL=info
```

Docker Compose will automatically use these variables.

#### Quick Commands (Using Makefile)

For convenience, use the Makefile:

```bash
# Show all available commands
make help

# Production
make build      # Build Docker images
make up         # Start all services
make down       # Stop all services
make logs       # View API logs
make clean      # Stop and remove volumes

# Development
make dev-up     # Start dev environment with hot reload
make dev-down   # Stop dev environment
make dev-logs   # View dev logs

# Database
make migrate-up   # Run migrations
make migrate-down # Rollback migrations
make db-shell     # Access PostgreSQL shell

# Testing
make test         # Run tests
make test-coverage # Run tests with coverage
```

## API Documentation

### Base URL

```
http://localhost:8080/api
```

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

### Endpoints

#### Auth Endpoints

**POST `/api/auth/register`** - Register new user
```json
{
  "email": "user@example.com",
  "password": "password123",
  "role": "client",  // or "provider"
  "full_name": "John Doe",
  "phone": "1234567890"
}
```

**POST `/api/auth/login`** - Login
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**GET `/api/auth/me`** - Get current user (protected)

#### Service Provider Endpoints

**GET `/api/providers`** - List providers (public)
Query parameters:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `service` - Filter by service ID (UUID)
- `city` - Filter by city name
- `min_rating` - Minimum rating (float)
- `available` - Filter by availability (true/false)
- `sort` - Sort field (rating, experience, hourly_rate)
- `order` - Sort order (asc, desc)

Example: `/api/providers?service=<uuid>&city=Almaty&min_rating=4&page=1&limit=10&sort=rating&order=desc`

**GET `/api/providers/:id`** - Get provider details (public)

**PUT `/api/providers/:id`** - Update provider profile (protected, own profile only)
```json
{
  "bio": "Experienced plumber",
  "hourly_rate": 25.50,
  "experience_years": 5,
  "location": "Almaty, Kazakhstan",
  "is_available": true
}
```

**DELETE `/api/providers/:id`** - Deactivate account (protected, own profile only)

#### Booking Endpoints

**POST `/api/bookings`** - Create booking (protected, client only)
```json
{
  "provider_id": "<uuid>",
  "service_id": "<uuid>",
  "scheduled_date": "2024-01-15T10:00:00Z",
  "duration_hours": 2.0,
  "notes": "Please bring tools"
}
```

**GET `/api/bookings`** - List user's bookings (protected)
- Returns bookings filtered by user role (client sees their bookings, provider sees bookings for them)

**GET `/api/bookings/:id`** - Get booking details (protected)

**PATCH `/api/bookings/:id/status`** - Update booking status (protected)
```json
{
  "status": "accepted"  // or "rejected", "completed", "cancelled"
}
```
- Providers can: accept, reject
- Clients can: cancel

#### Review Endpoints

**POST `/api/reviews`** - Create review (protected, client only)
```json
{
  "booking_id": "<uuid>",
  "rating": 5,  // 1-5
  "comment": "Great service!"
}
```

**GET `/api/reviews/providers/:id`** - Get provider's reviews (public)

#### Service Category Endpoints

**GET `/api/services`** - List all service categories (public)

**POST `/api/services`** - Create service category (protected, admin only)
```json
{
  "name": "Plumbing",
  "description": "Plumbing services",
  "icon_url": "https://example.com/icon.png"
}
```

**GET `/api/services/:id`** - Get service details (public)

### Response Format

Success response:
```json
{
  "data": {...},
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

Error response:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Too Many Requests (rate limit)
- `500` - Internal Server Error

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific test:

```bash
go test ./tests -run TestValidateEmail
```

## Database Schema

### Tables

1. **users** - User accounts (clients, providers, admins)
2. **service_providers** - Provider profiles (one-to-one with users)
3. **services** - Service categories (Plumbing, Electrical, etc.)
4. **provider_services** - Many-to-many relationship between providers and services
5. **bookings** - Service bookings/appointments
6. **reviews** - Reviews and ratings for completed bookings

### Relationships

- `users` 1:1 `service_providers`
- `service_providers` M:M `services` (via `provider_services`)
- `users` 1:M `bookings` (as client)
- `service_providers` 1:M `bookings`
- `bookings` 1:1 `reviews`

## Architecture

### Layers

1. **Handler Layer** - HTTP request/response handling, validation
2. **Service Layer** - Business logic, orchestration
3. **Repository Layer** - Database access, data persistence
4. **Model Layer** - Domain entities and DTOs

### Middleware

- **Authentication** - JWT token validation
- **Authorization** - Role-based access control
- **Logging** - Request logging with request IDs
- **Rate Limiting** - 100 requests per minute per IP
- **CORS** - Cross-origin resource sharing
- **Recovery** - Panic recovery

### Background Workers

- **Notification Worker** - Processes booking notifications asynchronously
  - Worker pool with 5 goroutines
  - Graceful shutdown support
  - Queue-based task processing

## Development

### Code Style

- Follow Go naming conventions
- Use interfaces for dependency injection
- Context propagation in all functions
- Structured error handling
- Use `defer` for cleanup

### Dependencies

Key dependencies:
- `github.com/gin-gonic/gin` - HTTP router
- `github.com/jmoiron/sqlx` - SQL extensions
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/golang-migrate/migrate/v4` - Database migrations

## Production Considerations

1. **Environment Variables**: Use secure secrets management (Docker secrets, Kubernetes secrets, etc.)
2. **Database**: Use connection pooling (configured in main.go)
3. **Logging**: Implement structured logging with proper log levels
4. **Rate Limiting**: Adjust limits based on traffic
5. **CORS**: Configure allowed origins for production
6. **HTTPS**: Use TLS in production (consider using a reverse proxy like nginx)
7. **Monitoring**: Add health checks and metrics
8. **Docker**: Use multi-stage builds for smaller image sizes (already implemented)
9. **Security**: 
   - Never commit `.env` files with real secrets
   - Use Docker secrets or environment variable injection
   - Regularly update base images for security patches

## License

MIT

## Author

Qasynda Development Team

