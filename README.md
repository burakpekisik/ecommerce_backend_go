# Ecommerce Backend (Go)

A backend service for an ecommerce application built with Go. This service handles essential ecommerce operations such as user management, product catalog, and order processing.

## Features

- **User Authentication**: Sign up, login, and session management.
- **Product Management**: Add, update, and remove products.
- **Order Processing**: Handle customer orders and track statuses.
- **Database Integration**: Uses a relational database for data storage.

## Getting Started

### Prerequisites

- Go 1.18+
- Database (MySQL/PostgreSQL)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/burakpekisik/ecommerce_backend_go.git
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the configuration file (`config/config.json`) with your database and authentication details.

4. Run the service:
   ```bash
   go run cmd/main.go
   ```

## Folder Structure

```plaintext
├── bin/             # Compiled binaries
├── cmd/             # Application entry point
├── config/          # Configuration files
├── db/              # Database migrations and setup
├── service/         # Core services (user, product, order)
├── types/           # Struct definitions
└── utils/           # Helper functions
```

## Database Setup

- Use the provided SQL scripts in `db/` to set up the necessary tables for users, products, and orders.

## License

This project is licensed under the MIT License.
