# Quick Start Guide

## 🚀 Fastest Way: Using Docker (Recommended)

### Step 1: Start Everything

```bash
# Navigate to project directory
cd "/Users/batyrkhan/Desktop/Golang dev/go-final"

# Start all services (PostgreSQL + API + Migrations)
docker-compose up -d
```

This will:
- ✅ Start PostgreSQL database
- ✅ Run database migrations automatically
- ✅ Build and start the API server

### Step 2: Check if it's Running

```bash
# View logs to see if everything started correctly
docker-compose logs -f api
```

You should see: `Server starting on port 8080`

### Step 3: Test the API

Open your browser or use curl:

```bash
# Health check
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### Step 4: Access the API

The API is now running at: **http://localhost:8080**

Try registering a user:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "role": "client",
    "full_name": "Test User",
    "phone": "1234567890"
  }'
```

---

## 🛠️ Development Mode (with Hot Reload)

For development with automatic code reloading:

```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up
```

Now any changes you make to `.go` files will automatically rebuild and restart the server!

---

## 📋 Useful Commands

### Using Makefile (Easier)

```bash
make help        # Show all commands
make up          # Start production
make down        # Stop services
make logs        # View logs
make dev-up      # Start development mode
make db-shell    # Access database
```

### Using Docker Compose Directly

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Stop and remove all data (fresh start)
docker-compose down -v

# Rebuild after code changes
docker-compose up -d --build
```

---

## 🏃 Local Development (Without Docker)

If you prefer to run locally:

### Step 1: Install Dependencies

```bash
go mod download
```

### Step 2: Set up PostgreSQL

Make sure PostgreSQL is running locally and create the database:

```bash
createdb qasynda
# OR
psql -U postgres -c "CREATE DATABASE qasynda;"
```

### Step 3: Run Migrations

```bash
# Install golang-migrate first
brew install golang-migrate  # macOS
# OR
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/qasynda?sslmode=disable" up
```

### Step 4: Create .env file

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### Step 5: Run the Server

```bash
go run cmd/api/main.go
```

---

## 🐛 Troubleshooting

### Port Already in Use

If port 8080 is already in use:

```bash
# Change port in docker-compose.yml or .env file
PORT=8081
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps

# View database logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Reset Everything

```bash
# Stop and remove everything
docker-compose down -v

# Start fresh
docker-compose up -d
```

---

## 📚 Next Steps

1. **Test the API**: Use Postman, curl, or your browser
2. **Read the README**: Full API documentation is in README.md
3. **Check endpoints**: Start with `/api/auth/register` and `/api/auth/login`

---

## 🎯 Quick Test

After starting, test these endpoints:

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"pass1234","role":"client","full_name":"Test","phone":"123"}'

# 3. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"pass1234"}'

# 4. Get services list
curl http://localhost:8080/api/services
```

Happy coding! 🚀

