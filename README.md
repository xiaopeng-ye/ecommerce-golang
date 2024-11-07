# Go E-commerce Backend Service

This is an e-commerce backend service developed using Go, providing a complete set of core e-commerce APIs, including user authentication, product management, shopping cart, and order system.

## Features

- User System
  - User registration
  - User login (JWT authentication)
  - User information management
- Product System
  - Product list
  - Product details
  - Product inventory management
- Shopping Cart System
  - Add products to shopping cart
  - Shopping cart checkout
- Order System
  - Order creation
  - Order management

## Tech Stack

- Go
- MySQL
- JWT Authentication
- Gorilla Mux (HTTP Router)
- SQL Migration Tool

## Project Structure

```bash
.
├── cmd/                    # Main application entry points
│   ├── api/               # API server configuration
│   ├── main.go            # Main program entry
│   └── migrate/           # Database migration tool
├── config/                # Configuration files
├── db/                    # Database connection and management
├── service/               # Business logic layer
│   ├── auth/             # Authentication service
│   ├── cart/             # Shopping cart service
│   ├── order/            # Order service
│   ├── product/          # Product service
│   └── user/             # User service
├── types/                 # Type definitions
└── utils/                 # Utility functions

```

## Environment Requirements

- Go 1.16+
- MySQL 8.0+
- Make

## Installation and Running

1. Clone the project

```bash
git clone https://github.com/xiaopeng-ye/ecommerce-golang.git
cd ecommerce-golang
```

2. Configure environment variables

```bash
cp .env.example .env
# Edit the .env file to configure database connection and other details
```

3. Run database migration

```bash
make migrate-up
```

4. Build and run the project

```bash
make run
```

## Available Make Commands

- `make build` - Build the project
- `make test` - Run tests
- `make run` - Build and run the project
- `make migration name=migration_name` - Create a new migration file
- `make migrate-up` - Run database migration
- `make migrate-down` - Rollback database migration

## API Endpoints

### User

- POST `/register` - User registration
- POST `/login` - User login
- GET `/users/{userID}` - Get user information (requires authentication)

### Product

- GET `/products` - Get product list
- GET `/products/{productID}` - Get product details
- POST `/products` - Add a new product (requires authentication)

### Order

- POST `/cart/checkout` - Shopping cart checkout (requires authentication)
- GET `/orders` - Get order list (requires authentication)

## Authentication

The API uses JWT (JSON Web Token) for authentication. In endpoints requiring authentication, include the following HTTP header:

```bash
Authorization: Bearer <your_jwt_token>
```

## Database Schema

The project uses a MySQL database, which primarily includes the following tables:

- users - User information
- products - Product information
- orders - Order information
- order_items - Order item details

## Testing

Run all tests:：

```bash
make test
```

## License

MIT License
