# Codingin - E-Commerce Platform untuk Digital Products

Platform marketplace untuk menjual source code, PDF, template, dan menerima pesanan custom development dengan sistem pembayaran manual (upload bukti transfer).

## ✨ Fitur Lengkap

### 🔐 Authentication & User Management

- ✅ Register/Login dengan Email & Password (bcrypt hashing)
- ✅ Google OAuth 2.0 Login
- ✅ GitHub OAuth Login
- ✅ JWT Token Authentication
- ✅ User Profile Management
- ✅ Change Password
- ✅ Admin User Management

### 🛍️ Product Management

- ✅ CRUD Categories (dengan slug auto-generate)
- ✅ CRUD Products (dengan image upload)
- ✅ Product Search & Filtering
- ✅ Category-based Products
- ✅ Featured Products
- ✅ Pagination Support

### 🛒 Shopping Experience

- ✅ Shopping Cart (Add, Update, Remove, Clear)
- ✅ Wishlist Management
- ✅ Price Calculation (dengan discount support)

### 💰 Order & Payment System

- ✅ **Checkout Order** - Buat pesanan dari product
- ✅ **Upload Bukti Transfer** - User upload payment proof
- ✅ **Admin Verification** - Approve/Reject payment manually
- ✅ **Order History** - Track semua orders
- ✅ **Order Status** - pending → processing → completed/cancelled
- ✅ **Transaction Tracking** - Payment records dengan metadata

### 📥 Download Management

- ✅ Download Tracking - Record every download
- ✅ Download History - User download history
- ✅ Access Control - Only paid orders can download
- ✅ Download Counter - Track total downloads per product

### ⭐ Reviews & Ratings

- ✅ **Product Reviews** - User dapat memberikan review setelah membeli
- ✅ **Rating System** - 1-5 stars rating
- ✅ **Average Rating** - Hitung rata-rata rating per product
- ✅ **Review Management** - User dapat edit/delete review sendiri
- ✅ **Admin Moderation** - Admin dapat hapus review yang tidak pantas

### 🎨 Custom Orders

- ✅ **Request Custom Development** - User request project custom
- ✅ **Budget & Requirements** - Tentukan budget dan requirements
- ✅ **Admin Quotation** - Admin kasih harga dan estimasi waktu
- ✅ **Status Tracking** - pending → reviewing → quoted → in_progress → completed
- ✅ **Custom Order Management** - Track semua custom order requests

### 🔔 Notifications

- ✅ **Real-time Notifications** - User dapat notifikasi penting
- ✅ **Notification Types** - Order, Payment, Download, Review, Custom Order, System
- ✅ **Read/Unread Status** - Mark notification sebagai dibaca
- ✅ **Notification History** - Lihat semua notifikasi dengan pagination
- ✅ **Delete Notifications** - Hapus notifikasi yang tidak diperlukan

### 📊 Admin Dashboard & Analytics

- ✅ **Dashboard Overview** - Total users, products, orders, revenue
- ✅ **Revenue Stats** - Revenue analytics dengan date range
- ✅ **Top Products** - Product terlaris
- ✅ **User Statistics** - User registrations, roles breakdown
- ✅ **Order Statistics** - Order by status, payment status, conversion rate

### 🔒 Security & Middleware

- ✅ JWT Authentication Middleware
- ✅ Admin Authorization Middleware
- ✅ Rate Limiting (100 req/min per IP)
- ✅ CORS Configuration
- ✅ Request Logger
- ✅ API Monitoring & Logging
- ✅ Error Handler

## 🏗️ Tech Stack

**Backend:**

- **Framework:** Gin (Go 1.24)
- **ORM:** GORM
- **Database:** PostgreSQL 15+
- **Authentication:** JWT (golang-jwt/jwt/v5)
- **OAuth:** golang.org/x/oauth2
- **Password Hash:** bcrypt

**Infrastructure:**

- **Container:** Docker & Docker Compose
- **File Storage:** Local filesystem (./uploads)

## 📁 Clean Architecture Structure

```
backend_codingin/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # App configuration
│   ├── handlers/                # HTTP Controllers
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── category_handler.go
│   │   ├── product_handler.go
│   │   ├── cart_handler.go
│   │   ├── wishlist_handler.go
│   │   ├── order_handler.go
│   │   └── download_handler.go
│   ├── services/                # Business Logic Layer
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── category_service.go
│   │   ├── product_service.go
│   │   ├── cart_service.go
│   │   ├── wishlist_service.go
│   │   ├── order_service.go
│   │   └── download_service.go
│   ├── repositories/            # Data Access Layer
│   │   ├── user_repository.go
│   │   ├── category_repository.go
│   │   ├── product_repository.go
│   │   ├── cart_repository.go
│   │   ├── wishlist_repository.go
│   │   ├── order_repository.go
│   │   ├── transaction_repository.go
│   │   ├── download_repository.go
│   │   └── api_log_repository.go
│   ├── models/                  # Database Models (13 models)
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── transaction.go
│   │   ├── download.go
│   │   ├── cart.go
│   │   ├── wishlist.go
│   │   ├── review.go
│   │   ├── custom_order.go
│   │   ├── notification.go
│   │   ├── analytics.go
│   │   └── api_log.go
│   ├── middleware/              # HTTP Middleware
│   │   ├── auth.go              # JWT verification
│   │   ├── admin.go             # Admin authorization
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── rate_limiter.go
│   │   └── api_monitor.go
│   └── routes/
│       └── routes.go            # API route definitions
├── pkg/
│   └── utils/                   # Shared utilities
│       ├── jwt.go               # JWT helpers
│       ├── oauth.go             # OAuth helpers
│       ├── response.go          # Response formatter
│       ├── database.go          # DB connection
│       ├── helpers.go
│       └── file_upload.go       # File upload handler
├── docs/                        # Documentation
│   ├── api-structure.md
│   └── database-schema.md
├── uploads/                     # File storage (auto-created)
│   ├── products/
│   └── payment_proofs/
├── .env                         # Environment config
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (optional)
- PostgreSQL 15+ (jika tidak pakai Docker)

### Installation

#### 1. Clone Repository

```bash
git clone <repository-url>
cd backend_codingin
```

#### 2. Setup Environment

```bash
cp .env.example .env
```

Edit `.env` dengan konfigurasi Anda:

```env
# App Configuration
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=codingin_db

# JWT
JWT_SECRET=your-super-secret-key-min-32-chars

# OAuth (optional)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback
```

#### 3. Run dengan Docker (Recommended)

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop
docker-compose down
```

#### 4. Run Manual (Tanpa Docker)

```bash
# Install dependencies
go mod download

# Create database
createdb codingin_db

# Run migration (auto on startup)
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Documentation

Base URL: `http://localhost:8080`

### 🔐 Authentication

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### OAuth Login

```http
GET /api/v1/auth/google
GET /api/v1/auth/github
```

### 👤 User Management (Protected)

```http
GET    /api/v1/user/profile          # Get profile
PUT    /api/v1/user/profile          # Update profile
PUT    /api/v1/user/password         # Change password
DELETE /api/v1/user/account          # Delete account
```

### 📦 Categories (Public Read, Admin Write)

```http
GET    /api/v1/categories            # Get all
GET    /api/v1/categories/:id        # Get by ID
GET    /api/v1/categories/slug/:slug # Get by slug
POST   /api/v1/categories            # Create (Admin)
PUT    /api/v1/categories/:id        # Update (Admin)
DELETE /api/v1/categories/:id        # Delete (Admin)
```

### 🛍️ Products (Public Read, Admin Write)

```http
GET    /api/v1/products                          # Get all + pagination
GET    /api/v1/products/featured                 # Featured products
GET    /api/v1/products/:id                      # Get by ID
GET    /api/v1/products/slug/:slug               # Get by slug
GET    /api/v1/products/category/:category_id    # By category
POST   /api/v1/products                          # Create (Admin)
PUT    /api/v1/products/:id                      # Update (Admin)
DELETE /api/v1/products/:id                      # Delete (Admin)
```

Query Parameters:

- `?page=1&limit=10` - Pagination
- `?category_id=1` - Filter by category
- `?search=keyword` - Search in title/description

### 🛒 Shopping Cart (Protected)

```http
POST   /api/v1/cart              # Add to cart
GET    /api/v1/cart              # Get cart + total
PUT    /api/v1/cart/:id          # Update quantity
DELETE /api/v1/cart/:id          # Remove item
DELETE /api/v1/cart/clear        # Clear cart
```

### ❤️ Wishlist (Protected)

```http
POST   /api/v1/wishlist          # Add to wishlist
GET    /api/v1/wishlist          # Get wishlist
DELETE /api/v1/wishlist/:id      # Remove item
DELETE /api/v1/wishlist/clear    # Clear wishlist
```

### 💰 Orders (Protected)

```http
POST   /api/v1/orders                      # Create order (checkout)
GET    /api/v1/orders                      # Get user orders
GET    /api/v1/orders/:id                  # Get order detail
POST   /api/v1/orders/:id/payment-proof    # Upload bukti transfer
POST   /api/v1/orders/:id/cancel           # Cancel order
```

**Create Order Request:**

```json
{
  "product_id": 1,
  "quantity": 1
}
```

**Upload Payment Proof:**

```http
POST /api/v1/orders/:id/payment-proof
Content-Type: multipart/form-data

proof: (image file)
```

### 📥 Downloads (Protected)

```http
POST /api/v1/downloads?product_id=1&order_id=1  # Download product
GET  /api/v1/downloads                          # Download history
GET  /api/v1/downloads/history/:product_id      # Product download history
```

### ⭐ Reviews (Public Read, Protected Write)

```http
GET    /api/v1/reviews/product/:product_id     # Get product reviews
POST   /api/v1/reviews                          # Create review (Protected)
GET    /api/v1/reviews/me                       # My reviews (Protected)
PUT    /api/v1/reviews/:id                      # Update review (Protected)
DELETE /api/v1/reviews/:id                      # Delete review (Protected)
```

**Create Review Request:**

```json
{
  "product_id": 1,
  "rating": 5,
  "comment": "Sangat membantu! Source code-nya rapih dan lengkap"
}
```

### 🎨 Custom Orders (Protected)

```http
POST   /api/v1/custom-orders               # Create custom order request
GET    /api/v1/custom-orders/me            # My custom orders
GET    /api/v1/custom-orders/:id           # Get custom order detail
PUT    /api/v1/custom-orders/:id/cancel    # Cancel custom order
```

**Create Custom Order Request:**

```json
{
  "title": "Website E-Commerce dengan Laravel",
  "description": "Butuh website toko online full stack dengan Laravel & Vue.js",
  "requirements": "- Laravel 10\n- Vue 3 + Vite\n- Payment Gateway\n- Admin Dashboard",
  "budget": 5000000
}
```

### 🔔 Notifications (Protected)

```http
GET    /api/v1/notifications                # Get notifications + pagination
GET    /api/v1/notifications/unread         # Get unread notifications
PUT    /api/v1/notifications/:id/read       # Mark as read
PUT    /api/v1/notifications/read-all       # Mark all as read
DELETE /api/v1/notifications/:id            # Delete notification
```

### 👨‍💼 Admin Endpoints

#### User Management

```http
GET    /api/v1/admin/users           # Get all users
GET    /api/v1/admin/users/:id       # Get user by ID
PUT    /api/v1/admin/users/:id       # Update user
DELETE /api/v1/admin/users/:id       # Delete user
```

#### Order Management

```http
GET  /api/v1/admin/orders                 # Get all orders
POST /api/v1/admin/orders/:id/approve     # Approve payment
POST /api/v1/admin/orders/:id/reject      # Reject payment
```

**Reject Payment Request:**

```json
{
  "reason": "Bukti transfer tidak jelas"
}
```

#### Custom Orders Management

```http
GET  /api/v1/admin/custom-orders              # Get all custom orders
PUT  /api/v1/admin/custom-orders/:id/process  # Process custom order
```

**Process Custom Order Request:**

```json
{
  "status": "quoted",
  "admin_notes": "Estimasi pengerjaan 2 minggu",
  "quoted_price": 7500000,
  "estimated_days": 14
}
```

Status options: `pending`, `reviewing`, `quoted`, `in_progress`, `completed`, `cancelled`

#### Reviews Moderation

```http
DELETE /api/v1/admin/reviews/:id          # Delete review
```

#### Analytics & Dashboard

```http
GET /api/v1/admin/analytics/dashboard      # Dashboard overview
GET /api/v1/admin/analytics/revenue        # Revenue statistics
GET /api/v1/admin/analytics/top-products   # Top selling products
GET /api/v1/admin/analytics/users          # User statistics
GET /api/v1/admin/analytics/orders         # Order statistics
}
```

## 💳 Order Flow (Manual Payment)

### 1️⃣ Customer Flow

```
1. Browse Products
   GET /api/v1/products

2. Add to Cart
   POST /api/v1/cart
   { "product_id": 1, "quantity": 1 }

3. Checkout (Create Order)
   POST /api/v1/orders
   { "product_id": 1, "quantity": 1 }

   Response: Order with status "pending"

4. Transfer to Bank Account
   (Manual - diluar sistem)

5. Upload Payment Proof
   POST /api/v1/orders/{order_id}/payment-proof
   FormData: proof (image)

   Status: pending → processing

6. Wait for Admin Verification
   (Check order status)

7. After Approved - Download Product
   POST /api/v1/downloads?product_id=1&order_id=1

   Status: processing → completed
```

### 2️⃣ Admin Flow

```
1. View Pending Orders
   GET /api/v1/admin/orders?status=processing

2. View Payment Proof
   Check: /uploads/payment_proofs/{filename}

3. Approve Payment
   POST /api/v1/admin/orders/{order_id}/approve

   OR

4. Reject Payment
   POST /api/v1/admin/orders/{order_id}/reject
   { "reason": "Invalid proof" }
```

## 📊 Order Status Flow

```
pending (just created)
    ↓
processing (payment proof uploaded)
    ↓
completed (approved by admin) → Customer can download
    OR
cancelled (rejected by admin / cancelled by user)
```

**Payment Status:**

- `pending` - Waiting for payment proof
- `paid` - Payment approved
- `failed` - Payment rejected
- `cancelled` - Order cancelled

## 🗄️ Database Models

**Core Models:**

- User (with role: user/admin)
- Product (digital products)
- Category
- Order (with quantity support)
- Transaction (payment tracking)
- Download (download tracking)
- Cart
- Wishlist
- Review (product reviews & ratings)
- CustomOrder (custom development requests)
- Notification (user notifications)
- Analytics (platform analytics)
- APILog (API request tracking)

## 🔒 Authentication & Authorization

**JWT Token:**

- Issued after successful login/register
- Include in header: `Authorization: Bearer {token}`
- Expires in 24 hours (configurable)

**Roles:**

- `user` - Regular customer
- `admin` - Full access to admin endpoints

**Protected Routes:**

- User routes: Requires valid JWT
- Admin routes: Requires JWT + admin role

## 📝 Response Format

**Success Response:**

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**

```json
{
  "success": false,
  "message": "Error message"
}
```

## 🧪 Testing

### Manual Testing with curl

**Register:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Login:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Get Profile (dengan token):**

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 Roadmap

**Completed ✅**

- [x] Authentication (JWT + OAuth Google & GitHub)
- [x] User Management (Profile, Password, Admin CRUD)
- [x] Product & Category CRUD dengan Image Upload
- [x] Shopping Cart & Wishlist
- [x] Order System dengan Manual Payment (Upload Bukti Transfer)
- [x] Download Tracking & Access Control
- [x] Admin Payment Verification System
- [x] Reviews & Ratings System
- [x] Custom Order Management (Request & Quotation)
- [x] Notifications System
- [x] Admin Dashboard & Analytics

**Planned 📋**

- [ ] Email Notifications (SMTP Integration)
- [ ] Payment Gateway Integration (DUITKU/Midtrans) - Optional
- [ ] Automatic Payment Verification
- [ ] Invoice Generation (PDF)
- [ ] Refund System
- [ ] Promo Code/Discount System
- [ ] Product Bundles
- [ ] Affiliate System
- [ ] Advanced Analytics (Charts & Reports)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Developer

**Codingin Team**

- Backend: Go + Gin Framework
- Database: PostgreSQL
- Architecture: Clean Architecture Pattern

---

**Happy Coding! 🚀**

{
"name": "John Doe",
"email": "john@example.com",
"password": "password123"
}

````

#### Get All Users

```bash
GET /api/v1/users
````

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
