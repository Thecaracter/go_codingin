# Codingin - Database Schema Design

## Entities

### 1. Users

- id (PK)
- email (unique)
- password (hashed, nullable for OAuth users)
- name
- role (user/admin)
- provider (local/google/github)
- provider_id
- avatar_url
- is_verified
- created_at, updated_at, deleted_at

### 2. Products (Source Code, PDF, etc)

- id (PK)
- title
- slug (unique)
- description
- category_id (FK)
- type (source_code/pdf/template/other)
- price
- discount_price
- preview_images (JSON array)
- demo_url
- file_url (encrypted/secured)
- tech_stack (JSON array)
- features (JSON array)
- requirements (JSON array)
- downloads_count
- views_count
- rating_average
- is_active
- created_by (FK to users)
- created_at, updated_at, deleted_at

### 3. Categories

- id (PK)
- name
- slug (unique)
- description
- icon
- parent_id (FK to categories - for sub-categories)
- order
- is_active
- created_at, updated_at

### 4. Orders

- id (PK)
- order_number (unique)
- user_id (FK)
- product_id (FK, nullable for custom orders)
- order_type (product/custom)
- status (pending/processing/completed/cancelled/refunded)
- total_amount
- discount_amount
- final_amount
- payment_method
- payment_status (pending/paid/failed/refunded)
- payment_id (from payment gateway)
- notes
- created_at, updated_at, deleted_at

### 5. Custom Orders

- id (PK)
- order_id (FK)
- user_id (FK)
- title
- description
- requirements (text)
- budget_min
- budget_max
- deadline
- status (submitted/under_review/in_progress/completed/cancelled)
- admin_notes (text)
- attachments (JSON array)
- quote_amount
- agreed_amount
- created_at, updated_at

### 6. Transactions

- id (PK)
- order_id (FK)
- user_id (FK)
- transaction_number (unique)
- amount
- payment_method (credit_card/bank_transfer/ewallet/crypto)
- payment_gateway (midtrans/stripe/xendit/etc)
- payment_gateway_ref
- status (pending/success/failed/cancelled)
- paid_at
- metadata (JSON)
- created_at, updated_at

### 7. Downloads

- id (PK)
- user_id (FK)
- product_id (FK)
- order_id (FK)
- download_url (temporary signed URL)
- expires_at
- is_used
- downloaded_at
- created_at

### 8. Reviews

- id (PK)
- user_id (FK)
- product_id (FK)
- order_id (FK)
- rating (1-5)
- comment
- is_verified_purchase
- created_at, updated_at, deleted_at

### 9. Carts

- id (PK)
- user_id (FK)
- product_id (FK)
- quantity (default 1 for digital products)
- created_at, updated_at

### 10. Wishlists

- id (PK)
- user_id (FK)
- product_id (FK)
- created_at

### 11. API Logs (for monitoring)

- id (PK)
- user_id (FK, nullable)
- method (GET/POST/PUT/DELETE)
- endpoint
- status_code
- response_time_ms
- ip_address
- user_agent
- request_body (JSON, truncated)
- response_body (JSON, truncated)
- error_message
- created_at

### 12. Analytics

- id (PK)
- date
- metric_type (page_view/unique_visitor/product_view/conversion/etc)
- metric_value (int)
- metadata (JSON - page, product_id, etc)
- created_at

### 13. Notifications

- id (PK)
- user_id (FK)
- type (order/payment/system/custom_order)
- title
- message
- action_url
- is_read
- read_at
- created_at

## Relationships

- Users 1:N Orders
- Users 1:N Custom Orders
- Users 1:N Reviews
- Users 1:N Downloads
- Users 1:N Carts
- Users 1:N Wishlists
- Products 1:N Orders
- Products 1:N Reviews
- Products 1:N Downloads
- Categories 1:N Products
- Categories 1:N Categories (self-referencing for sub-categories)
- Orders 1:1 Custom Orders
- Orders 1:N Transactions
- Orders 1:N Downloads

## Indexes

- users(email), users(provider, provider_id)
- products(slug), products(category_id), products(is_active)
- orders(user_id), orders(order_number), orders(status)
- custom_orders(user_id), custom_orders(status)
- transactions(order_id), transactions(user_id), transactions(status)
- reviews(product_id), reviews(user_id)
- api_logs(created_at), api_logs(endpoint)
- analytics(date), analytics(metric_type)
