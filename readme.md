# SpotSync API

SpotSync is a parking management RESTful API built with Go, Echo, and PostgreSQL. It provides authentication, parking zone management (for admins), and a reservation system for end users to book parking spots.

## 🚀 Features

- **User Authentication**: Registration, login, and JWT-based authentication. Role-based access control (`admin`, `user`).
- **Parking Zones Management**: Admins can create, update, delete, and list parking zones, including tracking total capacity and available spots.
- **Reservations**: Users can reserve parking spots by providing their license plate. Users can view their reservations and cancel active ones. Admins can view all reservations across the system.
- **Swagger Documentation**: Integrated Swagger UI for easy API testing and exploration.

## 🛠️ Tech Stack

- **Language**: [Go](https://go.dev/) (v1.26+)
- **Framework**: [Echo](https://echo.labstack.com/) (v4)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **Authentication**: JWT (`golang-jwt/jwt/v5`)
- **Validation**: `go-playground/validator/v10`
- **Documentation**: Swagger (`swaggo/echo-swagger`)

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

## 🔌 API Endpoints

### Health & Docs
- `GET /` - Health check
- `GET /swagger/index.html` - Swagger API Documentation

### Authentication (`/api/v1/auth`)
- `POST /register` - Register a new user
- `POST /login` - Authenticate and receive a JWT
- `GET /me` - Get current authenticated user details

### Parking Zones (`/api/v1/zones`)
- `GET /` - List all parking zones
- `GET /:id` - Get details of a specific parking zone
- `POST /` - Create a new parking zone (Admin only)
- `PUT /:id` - Update a parking zone (Admin only)
- `DELETE /:id` - Delete a parking zone (Admin only)

### Reservations (`/api/v1/reservations`)
- `POST /` - Create a new reservation
- `GET /my-reservations` - List all reservations for the logged-in user
- `GET /` - List all reservations in the system (Admin only)
- `DELETE /:id` - Cancel a reservation

## ⚙️ Setup and Installation

1. **Clone the repository:**
   ```bash
   git clone <repository_url>
   cd spotsync-api
   ```

2. **Configure Environment Variables:**
   Copy the example environment file and update it with your database credentials and JWT secret.
   ```bash
   cp .env.example .env
   ```

3. **Run the Application:**
   ```bash
   go run cmd/main.go
   ```
   The server will start (by default on port 8000, depending on the `.env` configuration).

## 🎥 References

- **Interview Video / Resources**: [Google Drive Link](https://drive.google.com/drive/folders/1ADN7_KCpD5SNs5nrM5crc01xsK-nNVs1?usp=drive_link)
