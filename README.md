# Qasynda

### Modern Service Marketplace Platform

A robust microservices-based platform connecting clients with professional service providers. Built with performance, scalability, and modern aesthetics in mind.

---

### 🚀 Stack

*   **Backend**: Go (Golang) 1.21+, REST (Gin)
*   **Infrastructure**: Docker, PostgreSQL, Redis, RabbitMQ

### 🏗 Architecture

The platform follows a **Microservices Architecture** communicating over **HTTP/JSON**:

1.  **Gateway Service** (`:8080`): API Gateway, Reverse Proxy & Rate Limiter.
2.  **User Service** (`:50051`): Authentication, Roles & Profiles.
3.  **Marketplace Service** (`:50052`): Services, Bookings, Categories.
4.  **Chat Service** (`:50053`): Real-time messaging with WebSocket & RabbitMQ.

### 📖 API Endpoints (Gateway `:8080`)

#### Public
- `POST /api/auth/register` - Create new user
- `POST /api/auth/login` - Get JWT token
- `GET /api/services` - List available services
- `GET /api/providers` - Get provider list

#### Protected (Requires Bearer Token)
- `GET /api/auth/me` - Get current user profile
- `POST /api/services` - Create new service (Provider only)
- `POST /api/bookings` - Book a service
- `GET /api/bookings` - List my bookings
- `PUT /api/bookings/:id/status` - Update booking status
- `PUT /api/providers/status` - Toggle availability
- `GET /api/chat/history` - Get message history
- `WS /ws?user_id=...` - Real-time chat connection

### ⚡️ Quick Start

You only need **Docker** and **Make**.

**1. Start Infrastructure**
```bash
make reset-db
```

**2. Run Services**
Open separate terminals for each service:

```bash
# Terminal 1: Gateway
make run-gateway
```

```bash
# Terminal 2: User Service
make run-user
```

```bash
# Terminal 3: Marketplace Service
make run-marketplace
```

```bash
# Terminal 4: Chat Service
make run-chat
```

### 📂 Directory Structure

```
qasynda/
├── services/        # Microservices (Go/Gin)
│   ├── gateway/     # Entry point & Rate Limiting
│   ├── user/        # Auth & Users
│   ├── marketplace/ # Core Marketplace Logic
│   └── chat/        # Messaging & RMQ
├── shared/          # Shared Models, Auth & Libs
└── migrations/      # SQL Migrations
```

---

_© 2025 Qasynda Team_
