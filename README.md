# Order Service API

A RESTful API for managing orders built with Go. This project demonstrates a production-ready microservice with two implementations: in-memory storage and PostgreSQL persistence.

## Features

- Create, read, update, and delete orders
- RESTful API design
- JSON request/response format
- Input validation
- Error handling
- Two implementations (in-memory and PostgreSQL)

## Tech Stack

- Go (standard library net/http)
- PostgreSQL
- github.com/lib/pq (PostgreSQL driver)

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/orders` | Get all orders |
| POST | `/orders` | Create a new order |
| GET | `/orders/{id}` | Get a specific order |
| PUT | `/orders/{id}` | Update an order |
| DELETE | `/orders/{id}` | Delete an order |

## Request Examples

### Create an Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"item":"laptop","quantity":1}'
```

### Get All Orders
```bash
curl http://localhost:8080/orders
```

### Update an Order
```bash
curl -X PUT http://localhost:8080/orders/1 \
  -H "Content-Type: application/json" \
  -d '{"item":"gaming laptop","quantity":2}'
```

### Delete an Order
```bash
curl -X DELETE http://localhost:8080/orders/1
```

## Branches
This repository contains two branches showing the evolution of the project:
- main: Simple implementation with in-memory storage using a Go slice. Perfect for understanding the basic API structure.
- add-postgresql: Production-ready version with PostgreSQL database integration. Data persists across server restarts.

## Getting Started
#### Prerequisites
- Go 1.16 or higher
- PostgreSQL (for the database branch)

### Running the In-Memory Version (main branch)

```bash
git checkout main
go run main.go
```

The server will start at http://localhost:8080

### Running the PostgreSQL Version (add-postgresql branch)
```bash
git checkout add-postgresql
```

1. Create a PostgreSQL database:

```sql
CREATE DATABASE orderdb;
```

2. Create the orders table:

```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    item VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
3. Run the application:

```bash
go run main.go
```

#### Project Structure
```text
.
├── main.go          # Main application code
├── go.mod           # Go module file
├── go.sum           # Go module checksum
└── README.md        # This file
```

#### Error Handling
    The API returns appropriate HTTP status codes:

    200 OK: Successful request

    201 Created: Resource created successfully

    204 No Content: Resource deleted successfully

    400 Bad Request: Invalid input

    404 Not Found: Resource not found

    405 Method Not Allowed: HTTP method not supported

    500 Internal Server Error: Server-side error
    
Author
```
NabeelOG
```
