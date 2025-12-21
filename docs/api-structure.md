# Codingin - API Structure

Backend Go + Frontend Vue

## Folder Structure (Diperluas)

```
backend_codingin/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── order.go
│   │   ├── custom_order.go
│   │   ├── transaction.go
│   │   ├── download.go
│   │   ├── review.go
│   │   ├── cart.go
│   │   ├── wishlist.go
│   │   ├── api_log.go
│   │   ├── analytics.go
│   │   └── notification.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── category_repository.go
│   │   ├── order_repository.go
│   │   ├── custom_order_repository.go
│   │   ├── transaction_repository.go
│   │   ├── download_repository.go
│   │   ├── review_repository.go
│   │   ├── cart_repository.go
│   │   ├── wishlist_repository.go
│   │   ├── api_log_repository.go
│   │   ├── analytics_repository.go
│   │   └── notification_repository.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── category_service.go
│   │   ├── order_service.go
│   │   ├── custom_order_service.go
│   │   ├── payment_service.go
│   │   ├── download_service.go
│   │   ├── review_service.go
│   │   ├── cart_service.go
│   │   ├── wishlist_service.go
│   │   ├── analytics_service.go
│   │   ├── notification_service.go
│   │   └── oauth_service.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── product_handler.go
│   │   ├── category_handler.go
│   │   ├── order_handler.go
│   │   ├── custom_order_handler.go
│   │   ├── payment_handler.go
│   │   ├── download_handler.go
│   │   ├── review_handler.go
│   │   ├── cart_handler.go
│   │   ├── wishlist_handler.go
│   │   ├── admin_handler.go
│   │   └── analytics_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── admin.go
│   │   ├── logger.go
│   │   ├── cors.go
│   │   ├── rate_limiter.go
│   │   └── api_monitor.go
│   └── routes/
│       ├── routes.go
│       ├── auth_routes.go
│       ├── user_routes.go
│       ├── admin_routes.go
│       └── public_routes.go
├── pkg/
│   └── utils/
│       ├── response.go
│       ├── database.go
│       ├── jwt.go
│       ├── oauth.go
│       ├── file_upload.go
│       ├── payment.go
│       └── email.go
└── migrations/
```

## API Endpoints

### Public APIs

#### Landing & Products

```
GET  /api/v1/products              - List products (with filters)
GET  /api/v1/products/:slug        - Product detail
GET  /api/v1/categories            - List categories
GET  /api/v1/categories/:slug      - Products by category
GET  /api/v1/products/:id/reviews  - Product reviews
```

#### Authentication

```
POST /api/v1/auth/register         - Register with email/password
POST /api/v1/auth/login            - Login
POST /api/v1/auth/google           - Google OAuth
POST /api/v1/auth/github           - GitHub OAuth
POST /api/v1/auth/refresh          - Refresh token
POST /api/v1/auth/forgot-password  - Forgot password
POST /api/v1/auth/reset-password   - Reset password
GET  /api/v1/auth/verify-email     - Verify email
```

### User APIs (Protected)

#### Profile

```
GET    /api/v1/user/profile        - Get profile
PUT    /api/v1/user/profile        - Update profile
PUT    /api/v1/user/password       - Change password
DELETE /api/v1/user/account        - Delete account
```

#### Shopping

```
GET    /api/v1/user/cart           - Get cart
POST   /api/v1/user/cart           - Add to cart
DELETE /api/v1/user/cart/:id       - Remove from cart

GET    /api/v1/user/wishlist       - Get wishlist
POST   /api/v1/user/wishlist       - Add to wishlist
DELETE /api/v1/user/wishlist/:id   - Remove from wishlist
```

#### Orders

```
GET  /api/v1/user/orders           - Get user orders
GET  /api/v1/user/orders/:id       - Order detail
POST /api/v1/user/orders           - Create order (checkout)
POST /api/v1/user/orders/:id/pay   - Process payment

GET  /api/v1/user/downloads        - Get purchased products
GET  /api/v1/user/downloads/:id    - Generate download link
```

#### Custom Orders

```
GET  /api/v1/user/custom-orders          - Get custom orders
GET  /api/v1/user/custom-orders/:id      - Custom order detail
POST /api/v1/user/custom-orders          - Create custom order
PUT  /api/v1/user/custom-orders/:id      - Update custom order
POST /api/v1/user/custom-orders/:id/pay  - Pay custom order
```

#### Reviews

```
POST /api/v1/user/reviews          - Add review
PUT  /api/v1/user/reviews/:id      - Update review
DELETE /api/v1/user/reviews/:id    - Delete review
```

#### Notifications

```
GET  /api/v1/user/notifications         - Get notifications
PUT  /api/v1/user/notifications/:id/read - Mark as read
DELETE /api/v1/user/notifications/:id    - Delete notification
```

### Admin APIs (Protected + Admin Role)

#### Dashboard & Analytics

```
GET /api/v1/admin/dashboard              - Dashboard stats
GET /api/v1/admin/analytics              - Analytics data
GET /api/v1/admin/analytics/sales        - Sales analytics
GET /api/v1/admin/analytics/visitors     - Visitor analytics
GET /api/v1/admin/analytics/performance  - API performance
```

#### User Management

```
GET    /api/v1/admin/users         - List users
GET    /api/v1/admin/users/:id     - User detail
PUT    /api/v1/admin/users/:id     - Update user
DELETE /api/v1/admin/users/:id     - Delete user
GET    /api/v1/admin/users/active  - Active users
```

#### Product Management

```
GET    /api/v1/admin/products         - List products
GET    /api/v1/admin/products/:id     - Product detail
POST   /api/v1/admin/products         - Create product
PUT    /api/v1/admin/products/:id     - Update product
DELETE /api/v1/admin/products/:id     - Delete product
POST   /api/v1/admin/products/:id/upload - Upload files
```

#### Category Management

```
GET    /api/v1/admin/categories      - List categories
POST   /api/v1/admin/categories      - Create category
PUT    /api/v1/admin/categories/:id  - Update category
DELETE /api/v1/admin/categories/:id  - Delete category
```

#### Order Management

```
GET  /api/v1/admin/orders            - List orders
GET  /api/v1/admin/orders/:id        - Order detail
PUT  /api/v1/admin/orders/:id/status - Update order status
GET  /api/v1/admin/orders/statistics - Order statistics
```

#### Custom Order Management

```
GET  /api/v1/admin/custom-orders           - List custom orders
GET  /api/v1/admin/custom-orders/:id       - Custom order detail
PUT  /api/v1/admin/custom-orders/:id/status - Update status
PUT  /api/v1/admin/custom-orders/:id/quote  - Send quote
```

#### Payment & Transaction

```
GET /api/v1/admin/transactions          - List transactions
GET /api/v1/admin/transactions/:id      - Transaction detail
GET /api/v1/admin/transactions/summary  - Transaction summary
POST /api/v1/admin/transactions/:id/refund - Process refund
```

#### API Monitoring

```
GET /api/v1/admin/api-logs              - API logs
GET /api/v1/admin/api-logs/stats        - API statistics
GET /api/v1/admin/api-logs/performance  - Performance metrics
```

## Payment Integration Options

1. **Midtrans** (Indonesia)
2. **Stripe** (International)
3. **Xendit** (Indonesia + SEA)
4. **PayPal**

## OAuth Providers

1. **Google OAuth 2.0**
2. **GitHub OAuth**

## File Storage Options

1. **Local Storage** (Development)
2. **AWS S3** (Production)
3. **Google Cloud Storage**
4. **Cloudflare R2**

## Tech Stack

### Backend

- Go + Gin
- GORM + PostgreSQL
- JWT Authentication
- OAuth 2.0
- Payment Gateway SDK
- Redis (for caching & rate limiting)

### Frontend (Vue)

- Vue 3 + Composition API
- Vue Router
- Pinia (State Management)
- TailwindCSS
- Axios
