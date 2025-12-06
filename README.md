# Qasynda

### Modern Service Marketplace Platform

A robust microservices-based platform connecting clients with professional service providers. Built with performance, scalability, and modern aesthetics in mind.

---

### 🚀 Stack

*   **Backend**: Go (Golang) 1.21+, gRPC, Gin
*   **Infrastructure**: Docker, PostgreSQL, Redis, RabbitMQ
*   **Frontend**: Vanilla JS, Modern CSS (Glassmorphism), RWD

### 🏗 Architecture

The platform follows a **Microservices Architecture**:

1.  **Gateway Service** (`:8080`): API Gateway & Static File Server.
2.  **User Service** (`:50051`): Authentication & Profiles.
3.  **Marketplace Service** (`:50052`): Services, Bookings, Providers.
4.  **Chat Service** (`:50053`): Real-time messaging (In Progress).

### ⚡️ Quick Start

You only need **Docker** and **Make**.

**1. Start Infrastructure**
```bash
make reset-db
```

**2. Run Services**
Open separate terminals for each service:

```bash
# Terminal 1: Gateway (Frontend + API)
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

**3. Explore**
Open **[http://localhost:8080](http://localhost:8080)** in your browser.

---

### 🛠 Development Commands

| Command | Description |
| :--- | :--- |
| `make docker-up` | Start generic infrastructure (Postgres, Redis, RabbitMQ) |
| `make reset-db` | **Wipe** database volume & restart (Fixes connection issues) |
| `make proto` | Regenerate gRPC protobuf code |
| `make tidy` | Clean up Go modules |

### 📂 Directory Structure

```
qasynda/
├── frontend/        # Modern Web UI (HTML/CSS/JS)
├── services/        # Microservices (Go)
│   ├── gateway/     # REST API & File Server
│   ├── user/        # gRPC User Service
│   ├── marketplace/ # gRPC Marketplace Service
│   └── chat/        # gRPC Chat Service
├── shared/          # Shared Proto & Config
└── migrations/      # SQL Migrations
```

---

_© 2025 Qasynda Team_
