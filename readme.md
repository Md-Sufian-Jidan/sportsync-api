# 🚗 SpotSync API

> **Smart Parking & EV Charging Reservation Platform**
>
> A centralized RESTful API for airports and malls to manage parking zones and handle high-demand EV charging spot reservations — built with Go, Echo, and PostgreSQL.

---

## 🔗 Quick Links

| Resource | Link |
|---|---|
| 📦 GitHub Repository | [github.com/Md-Sufian-Jidan/sportsync-api](https://github.com/Md-Sufian-Jidan/sportsync-api) |
| 🌐 Live Deployment | [spotsync-api-eem6.onrender.com](https://sportsync-api-eem6.onrender.com) |
| 🎥 Interview Video | [Google Drive](https://drive.google.com/drive/folders/1ADN7_KCpD5SNs5nrM5crc01xsK-nNVs1?usp=drive_link) |
| 📖 API Docs (Swagger) | `<live-url>/swagger/index.html` |

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Database Schema](#-database-schema)
- [API Endpoints](#-api-endpoints)
- [Local Setup](#-local-setup)
- [Environment Variables](#-environment-variables)
- [Deployment](#-deployment)

---

## 🧭 Overview

SpotSync solves a real-world challenge: **concurrent EV charging spot reservation**. When two drivers simultaneously try to book the last available spot, a naive implementation creates a race condition — both get confirmed, exceeding capacity. SpotSync uses **PostgreSQL row-level locking (`FOR UPDATE`) inside GORM transactions** to prevent this atomically.

The system supports two roles:
- **Driver** — Browse zones, make reservations, view & cancel their own bookings.
- **Admin** — Full driver permissions plus zone management and system-wide reservation visibility.

---

## ✨ Features

- **JWT Authentication** — Stateless auth with signed tokens; passwords hashed using bcrypt (cost 10–12).
- **Role-Based Access Control** — Middleware enforces `driver` vs `admin` permissions per endpoint.
- **Parking Zone Management** — Admins create/update/delete zones with type, capacity, and pricing.
- **Dynamic Availability** — `available_spots` computed in real-time (`total_capacity` minus active reservations).
- **Race-Condition-Safe Reservations** — GORM transactions + `SELECT ... FOR UPDATE` prevent overbooking under concurrent load.
- **Reservation Lifecycle** — Statuses: `active` → `cancelled` / `completed`. Cancellations free up the spot immediately.
- **Swagger UI** — Interactive API documentation available at `/swagger/index.html`.
- **Centralized Error Handling** — Consistent JSON error responses; raw ORM errors never leak to the client.

---

## 🛠️ Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/) | 1.22+ | Core language |
| [Echo](https://echo.labstack.com/) | v4 | HTTP framework |
| [GORM](https://gorm.io/) | latest | ORM & query builder |
| [PostgreSQL](https://www.postgresql.org/) | 15+ | Relational database (NeonDB / Supabase) |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt) | v5 | JWT generation & verification |
| [go-playground/validator](https://github.com/go-playground/validator) | v10 | DTO struct validation |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | latest | bcrypt password hashing |
| [swaggo/echo-swagger](https://github.com/swaggo/echo-swagger) | latest | Swagger UI integration |

---

## 🏛️ Architecture

SpotSync follows **Clean Architecture** with strict layer separation. No layer talks to a lower layer it doesn't own, and no layer skips levels.

```
HTTP Request
     │
     ▼
┌──────────────┐
│   Handler    │  ← Binds & validates DTO, extracts JWT claims,
│  (handler/)  │    calls Service, returns JSON response
└──────┬───────┘
       │ calls
       ▼
┌──────────────┐
│   Service    │  ← Business logic: hashes passwords, generates JWTs,
│  (service/)  │    enforces capacity rules, orchestrates Repositories
└──────┬───────┘
       │ calls
       ▼
┌──────────────┐
│  Repository  │  ← All database operations via GORM
│ (repository/)│    Handles transactions & row-level locks
└──────┬───────┘
       │ reads/writes
       ▼
┌──────────────┐
│  PostgreSQL  │
└──────────────┘
```

**Supporting layers:**

| Layer | Directory | Role |
|---|---|---|
| **DTO** | `dto/` | Request payloads & response shapes. GORM models are never exposed directly to the API. |
| **Models** | `models/` | GORM structs mapped to database tables. |
| **Middleware** | `internal/middleware/` | JWT verification, role enforcement (`admin`/`driver`). |
| **Config** | `internal/config/` | Env loading, DB connection, connection pool settings. |

**Dependency wiring in `main.go` (manual, no DI framework):**

```
Repository  →  Service  →  Handler  →  Echo Router
```

### Concurrency Safety — The EV Spot Bottleneck

When two drivers race to claim the last spot:

```
Driver A ─────┐
              ├──▶ BEGIN TRANSACTION
Driver B ─────┘         │
                        ▼
              SELECT * FROM parking_zones
              WHERE id = ? FOR UPDATE   ← Only ONE transaction holds the lock
                        │
              COUNT active reservations
                        │
              IF count < capacity → INSERT reservation
              ELSE → return ErrZoneFull (409 Conflict)
                        │
                        ▼
                   COMMIT / ROLLBACK
```

The `FOR UPDATE` row lock means the second transaction blocks until the first commits, then re-reads the actual count — preventing overbooking.

---

## 📁 Project Structure

This project follows a Domain-Driven Design (DDD) approach.

```text
├── cmd
│   └── main.go                  # Application entry point
├── docs                         # Swagger generated documentation
├── internal
│   ├── auth                     # JWT generation and validation
│   ├── config                   # Environment loading and DB connection
│   ├── domain                   # Core business logic divided into domains
│   │   ├── admin                # ParkingZone entity, handler, service, repo
│   │   ├── reservations         # Reservation entity, handler, service, repo
│   │   └── user                 # User entity, handler, service, repo
│   ├── httpResponse             # Standardized HTTP response formats
│   ├── middleware               # Echo middlewares (Auth, Role checking)
│   └── server                   # HTTP server and route registration
└── readme.md                    # Project documentation
```

---

## 🗄️ Database Schema

### `users`

| Column | Type | Constraints |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `name` | VARCHAR | NOT NULL |
| `email` | VARCHAR | NOT NULL, UNIQUE |
| `password` | VARCHAR | NOT NULL (bcrypt hash) |
| `role` | VARCHAR | NOT NULL, DEFAULT `'driver'`, CHECK IN (`driver`, `admin`) |
| `created_at` | TIMESTAMP | Auto-managed by GORM |
| `updated_at` | TIMESTAMP | Auto-managed by GORM |

### `parking_zones`

| Column | Type | Constraints |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `name` | VARCHAR | NOT NULL |
| `type` | VARCHAR | NOT NULL, CHECK IN (`general`, `ev_charging`, `covered`) |
| `total_capacity` | INTEGER | NOT NULL, > 0 |
| `price_per_hour` | DECIMAL | NOT NULL, > 0 |
| `created_at` | TIMESTAMP | Auto-managed by GORM |
| `updated_at` | TIMESTAMP | Auto-managed by GORM |

### `reservations`

| Column | Type | Constraints |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `user_id` | INTEGER | FOREIGN KEY → `users.id` |
| `zone_id` | INTEGER | FOREIGN KEY → `parking_zones.id` |
| `license_plate` | VARCHAR(15) | NOT NULL |
| `status` | VARCHAR | NOT NULL, DEFAULT `'active'`, CHECK IN (`active`, `completed`, `cancelled`) |
| `created_at` | TIMESTAMP | Auto-managed by GORM |
| `updated_at` | TIMESTAMP | Auto-managed by GORM |

---

## 🌐 API Endpoints

### Authentication — `/api/v1/auth`

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/register` | Public | Register a new user |
| `POST` | `/login` | Public | Login and receive a JWT |
| `GET` | `/me` | Authenticated | Get current user profile |

**Register** — `POST /api/v1/auth/register`
```json
// Request
{ "name": "John Doe", "email": "john@spotsync.com", "password": "securePass123", "role": "driver" }

// Response 201
{ "success": true, "message": "User registered successfully", "data": { "id": 1, "name": "John Doe", "email": "john@spotsync.com", "role": "driver", "created_at": "...", "updated_at": "..." } }
```

**Login** — `POST /api/v1/auth/login`
```json
// Request
{ "email": "john@spotsync.com", "password": "securePass123" }

// Response 200
{ "success": true, "message": "Login successful", "data": { "token": "eyJhbGci...", "user": { "id": 1, "name": "John Doe", "email": "john@spotsync.com", "role": "driver" } } }
```

---

### Parking Zones — `/api/v1/zones`

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/` | Public | List all zones with live `available_spots` |
| `GET` | `/:id` | Public | Get a single zone by ID |
| `POST` | `/` | Admin only | Create a new parking zone |
| `PUT` | `/:id` | Admin only | Update zone details |
| `DELETE` | `/:id` | Admin only | Delete a zone |

**Create Zone** — `POST /api/v1/zones` *(Admin only)*
```json
// Request
{ "name": "Terminal 1 EV Charging", "type": "ev_charging", "total_capacity": 20, "price_per_hour": 5.50 }

// Response 201
{ "success": true, "message": "Parking zone created successfully", "data": { "id": 5, "name": "Terminal 1 EV Charging", "type": "ev_charging", "total_capacity": 20, "price_per_hour": 5.50, "created_at": "...", "updated_at": "..." } }
```

**Get All Zones** — `GET /api/v1/zones`
```json
// Response 200
{ "success": true, "message": "Parking zones retrieved successfully", "data": [ { "id": 5, "name": "Terminal 1 EV Charging", "type": "ev_charging", "total_capacity": 20, "available_spots": 14, "price_per_hour": 5.50, "created_at": "..." } ] }
```

> `available_spots` is computed dynamically: `total_capacity − COUNT(active reservations)`.

---

### Reservations — `/api/v1/reservations`

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/` | Authenticated | Reserve a parking spot |
| `GET` | `/my-reservations` | Authenticated | View caller's own reservations |
| `GET` | `/` | Admin only | View all reservations system-wide |
| `DELETE` | `/:id` | Authenticated | Cancel a reservation (own only) |

**Reserve Spot** — `POST /api/v1/reservations`
```json
// Request
{ "zone_id": 5, "license_plate": "ABC-1234" }

// Response 201
{ "success": true, "message": "Reservation confirmed successfully", "data": { "id": 105, "user_id": 1, "zone_id": 5, "license_plate": "ABC-1234", "status": "active", "created_at": "...", "updated_at": "..." } }
```

**My Reservations** — `GET /api/v1/reservations/my-reservations`
```json
// Response 200
{ "success": true, "message": "My reservations retrieved successfully", "data": [ { "id": 105, "license_plate": "ABC-1234", "status": "active", "zone": { "id": 5, "name": "Terminal 1 EV Charging", "type": "ev_charging" }, "created_at": "..." } ] }
```

---

### HTTP Status Code Reference

| Code | Meaning | When Used |
|---|---|---|
| `200` | OK | Successful GET / DELETE |
| `201` | Created | Successful POST (resource created) |
| `400` | Bad Request | Validation errors, invalid input |
| `401` | Unauthorized | Missing, expired, or invalid JWT |
| `403` | Forbidden | Valid token, insufficient role/permissions |
| `404` | Not Found | Resource does not exist |
| `409` | Conflict | Zone full, duplicate license plate |
| `500` | Internal Server Error | Unexpected server/database error |

**Error response format:**
```json
{ "success": false, "message": "Error description", "errors": "Error details" }
```

---

## ⚙️ Local Setup

### Prerequisites

- Go 1.22 or higher installed ([download](https://go.dev/dl/))
- A running PostgreSQL instance (local, [NeonDB](https://neon.tech), or [Supabase](https://supabase.com))
- `git` installed

### Steps

**1. Clone the repository**
```bash
git clone https://github.com/Md-Sufian-Jidan/sportsync-api.git
cd sportsync-api
```

**2. Install dependencies**
```bash
go mod download
```

**3. Configure environment variables**
```bash
cp .env.example .env
# Open .env and fill in your values (see table below)
```

**4. Run the server**
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8000`.

**5. (Optional) Hot-reload during development**

Install [Air](https://github.com/air-verse/air) for live reloading:
```bash
go install github.com/air-verse/air@latest
air
```

---

## 🔑 Environment Variables

Create a `.env` file at the project root. All variables below are required.

| Variable | Example | Description |
|---|---|---|
| `PORT` | `8000` | Port the Echo server listens on |
| `Dsn` | `ep-xxx.us-east-2.aws.neon.tech` | PostgreSQL host |
| `JWT_SECRET` | `your_super_secret_key` | Secret key used to sign JWTs (keep this long and random) |


**`.env.example`**
```env
PORT=8000

Dsn=localhost

JWT_SECRET=change_me_to_a_long_random_string
```

---

## 🚀 Deployment

The API is deployed on **[Render](https://render.com)** with a **[NeonDB](https://neon.tech)** PostgreSQL database.

**Deployment checklist used:**
- Set all environment variables in Render's dashboard under *Environment*.
- `DB_SSLMODE=require` for NeonDB connections.
- CORS configured to allow the required origins via Echo's CORS middleware.
- Render auto-deploys on every push to the `main` branch.

**Build command (Render):**
```
go build -o spotsync ./cmd/main.go
```

**Start command (Render):**
```
./spotsync
```

---

## 🎥 Interview Video

Technical questions answered (any 3 of 5):

📎 [Watch on Google Drive](https://drive.google.com/drive/folders/1ADN7_KCpD5SNs5nrM5crc01xsK-nNVs1?usp=drive_link)

Topics covered include Go's concurrency model, GORM transactions with row-level locking, and interface duck typing.

---

## 📄 License

This project was built as part of a backend development assignment. All code is original work.