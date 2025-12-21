# Codingin - E-Commerce Platform untuk Source Code & Custom Development

Platform marketplace untuk menjual source code, PDF, template, dan menerima pesanan custom development.

## 🚀 Fitur

### User Features

- ✅ **Multi-Authentication**
  - Register/Login dengan Email & Password
  - Google OAuth Login
  - GitHub OAuth Login
  - JWT Token Authentication
- Docker & Docker Compose support
- Logging middleware
- CORS middleware
- Environment configuration
- Error handling

## 📁 Struktur Project

```
.
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Konfigurasi aplikasi
│   ├── handlers/
│   │   └── user_handler.go      # HTTP handlers
│   ├── models/
│   │   └── user.go              # Data models
│   ├── repositories/
│   │   └── user_repository.go   # Database layer
│   ├── services/
│   │   └── user_service.go      # Business logic
│   ├── middleware/
│   │   ├── logger.go            # Logging middleware
│   │   └── cors.go              # CORS middleware
│   └── routes/
│       └── routes.go            # Route definitions
├── pkg/
│   └── utils/
│       ├── response.go          # Response utilities
│       └── database.go          # Database utilities
├── .env.example                 # Environment template
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose configuration
└── go.mod
```

## 🛠️ Setup

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (untuk development dengan Docker)
- PostgreSQL (jika tidak pakai Docker)

### Cara Run

#### 1. Dengan Docker (Recommended)

```bash
# Copy file .env
cp .env.example .env

# Build dan jalankan semua services
docker-compose up -d

# Lihat logs
docker-compose logs -f app

# Stop services
docker-compose down
```

#### 2. Tanpa Docker

```bash
# Install dependencies
go mod download

# Copy file .env
cp .env.example .env

# Edit .env sesuai database lokal Anda
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=gin_db

# Run aplikasi
go run cmd/api/main.go
```

## 📡 API Endpoints

Base URL: `http://localhost:8080`

### Health Check

```
GET /health
```

### Users

#### Create User

```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Get All Users

```bash
GET /api/v1/users
```

#### Get User by ID

```bash
GET /api/v1/users/:id
```

#### Update User

```bash
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

#### Delete User

```bash
DELETE /api/v1/users/:id
```

## 🗄️ Database

PostgreSQL database akan otomatis dibuat oleh Docker Compose.

### Akses pgAdmin

- URL: `http://localhost:5050`
- Email: `admin@admin.com`
- Password: `admin`

### Connection Info

- Host: `postgres` (atau `localhost` jika di luar Docker)
- Port: `5432`
- Database: `gin_db`
- User: `postgres`
- Password: `postgres`

## 📦 Dependencies

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/sirupsen/logrus
```

## 🔧 Development

### Hot Reload (Optional)

Install Air untuk hot reload:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Migration

GORM Auto Migrate sudah dijalankan otomatis saat aplikasi start.

## 📝 Environment Variables

```env
APP_NAME=gin-quickstart
APP_ENV=development
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_db

JWT_SECRET=your-secret-key-change-this
```

## 🧪 Testing

Contoh testing dengan curl:

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Get all users
curl http://localhost:8080/api/v1/users

# Get user by ID
curl http://localhost:8080/api/v1/users/1

# Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated User","email":"updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 📄 License

MIT
