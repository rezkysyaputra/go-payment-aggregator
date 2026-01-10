# Go Payment Aggregator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust, enterprise-grade Payment Aggregator service designed to unify multiple payment gateways (Midtrans, Xendit, Stripe) under a single, standardized API. Built with **Go (Golang)**, adhering to **Clean Architecture** principles to ensure scalability, maintainability, and testability.

---

## 📑 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [Usage](#-usage)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Project Structure](#-project-structure)
- [License](#-license)

---

## 🚀 Features

- **Multi-Gateway Support**: Seamless integration with **Midtrans**, **Xendit**, and **Stripe** (planned).
- **Unified API Interface**: A single `CreateTransaction` endpoint intelligently routes requests to the appropriate provider.
- **Standardized Payment Methods**: Gateway-agnostic payment method format (`credit_card`, `bank_transfer`, `e_wallet`, `qris`).
- **Dynamic Gateway Selection**: Merchants can choose their preferred payment gateway per transaction.
- **Transaction Status Tracking**: Real-time transaction status checking across all gateways.
- **Resilient Webhook Handling**: Standardized webhook processing for payment notifications.
- **Merchant Callbacks**: Automatic notification system that relays payment status changes back to the merchant's registered `callback_url`.
- **Containerized**: Fully dockerized environment with PostgreSQL and Redis support for easy deployment.
- **Observability**: Structured logging with Logrus.
- **High Test Coverage**: 100% unit test coverage for business logic.

## 🏗 Architecture

This project follows **Hexagonal Architecture (Ports and Adapters)** to decouple business logic from external concerns.

```mermaid
graph TD
    Client["Client / Merchant"] -->|HTTP Request| Handler
    Handler -->|Call| Usecase["Usecase / Service"]
    Usecase -->|Persist| Repo["Repository (PostgreSQL)"]
    Usecase -->|Request Payment| Gateway["Payment Gateway Interface"]
    Gateway -->|Impl| Midtrans["Midtrans Adapter"]
    Gateway -->|Impl| Xendit["Xendit Adapter"]
    Gateway -->|Impl| Stripe["Stripe Adapter (Planned)"]
```

## 🛠 Tech Stack

- **Core**: Go 1.25+
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL
- **Caching**: Redis (Supported in config)
- **ORM**: GORM
- **Configuration**: Viper
- **Logging**: Logrus
- **Testing**: Testify (Suite, Assert, Mock)
- **Containerization**: Docker & Docker Compose

## 💻 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.25 or higher
- [Docker](https://www.docker.com/products/docker-desktop) & Docker Compose

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/rezkysyaputra/go-payment-aggregator.git
    cd go-payment-aggregator
    ```

2.  **Setup Configuration**
    Copy the example configuration file.
    ```bash
    cp .env.example .env
    ```

### Configuration

Edit `.env` with your credentials.

| Key | Description | Default |
| :--- | :--- | :--- |
| `APP_NAME` | Application Name | `go-payment-aggregator` |
| `SERVER_PORT` | HTTP Port | `8080` |
| `DATABASE_*` | Database connection details | - |
| `REDIS_*` | Redis connection details | - |
| `MIDTRANS_SERVER_KEY` | Midtrans Server Key | - |
| `MIDTRANS_ENVIRONMENT` | Midtrans Environment (`sandbox` or `production`) | `sandbox` |
| `XENDIT_API_KEY` | Xendit API Key | - |
| `STRIPE_SECRET_KEY` | Stripe Secret Key (optional, for future use) | - |
| `CONTEXT_TIMEOUT` | Request timeout in seconds | `2` |

## 🚀 Usage

### Running with Docker (Recommended)

The easiest way to run the application is using Docker Compose. This will start the API server and the PostgreSQL database.

```bash
docker-compose up -d --build
```

The API will be accessible at `http://localhost:8080`.

### Running Locally

1.  **Start Dependencies (DB)**
    ```bash
    docker-compose up -d db
    ```

2.  **Run Migrations**
    ```bash
    # Ensure migrations are applied (adjust command based on your setup)
    # Example: migrate -path database/migrations -database "postgresql://..." up
    ```

3.  **Start the Server**
    ```bash
    go run cmd/server/main.go
    ```

## 📡 API Documentation

The API follows RESTful conventions. A full OpenAPI specification is provided in `api/openapi.json`.

### Key Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/merchants` | Register a new merchant to get an API Key. |
| `GET` | `/api/v1/merchants/profile` | Get merchant profile (requires authentication). |
| `PUT` | `/api/v1/merchants/profile` | Update merchant profile. |
| `POST` | `/api/v1/merchants/api-key/regenerate` | Regenerate merchant API Key. |
| `POST` | `/api/v1/transactions` | Create a new transaction (supports `midtrans`, `xendit`). |
| `GET` | `/api/v1/transactions/{id}` | Retrieve transaction status by System ID. |
| `POST` | `/api/v1/webhooks/midtrans` | Webhook endpoint for Midtrans. |

### Standardized Payment Methods

The API uses gateway-agnostic payment method identifiers:

| Payment Method | Description | Supported Gateways |
| :--- | :--- | :--- |
| `credit_card` | Credit/Debit Card | Midtrans, Xendit, Stripe (planned) |
| `bank_transfer` | Bank Transfer / Virtual Account | Midtrans, Xendit |
| `e_wallet` | E-Wallet (GoPay, OVO, DANA, ShopeePay) | Midtrans, Xendit |
| `qris` | QR Code Payment | Midtrans, Xendit |

### Example: Create Transaction

```json
POST /api/v1/transactions
Content-Type: application/json
X-API-Key: mch_your_api_key_here

{
  "order_id": "ORDER-123456",
  "amount": 100000,
  "currency": "IDR",
  "provider": "xendit",
  "payment_method": "bank_transfer",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "items": [
    {
      "name": "Product A",
      "quantity": 2,
      "price": 50000
    }
  ]
}
```

**Response:**
```json
{
  "code": 201,
  "status": "success",
  "message": "Transaction created successfully",
  "data": {
    "id": "019b9836-deb7-7441-aaaf-ef51e4d91961",
    "order_id": "ORDER-123456",
    "amount": 100000,
    "currency": "IDR",
    "status": "PENDING",
    "payment_url": "https://checkout.xendit.co/web/...",
    "external_id": "62fe7ac7ae8faa001e3a7f01"
  }
}
```

## 🧪 Testing

This project includes both **Unit Tests** (for business logic) and **Integration Tests** (for end-to-end flows).

### Run All Tests
```bash
go test ./... -v
```

### Run Unit Tests Only
Unit tests are located within the `internal` packages (e.g., usecases).
```bash
go test ./internal/... -v
```

### Run Integration Tests Only
Integration tests are located in the `test` directory and require a running database.
```bash
go test ./test/... -v
```

## 📂 Project Structure

```
go-payment-aggregator/
├── api/                # OpenAPI/Swagger definitions
├── cmd/                # Main applications of the project
│   ├── server/         # API Server entrypoint
│   └── worker/         # Background worker entrypoint
├── internal/
│   ├── config/         # Configuration loading
│   ├── delivery/       # HTTP Handlers, Middleware, Routes
│   ├── domain/         # Business entities and interfaces (Core)
│   ├── gateway/        # 3rd Party API Adapters (Midtrans, Xendit, Stripe)
│   ├── mocks/          # Mocks generated by Mockery
│   ├── pkg/            # Internal shared packages (Crypto, UUID, etc.)
│   ├── repository/     # Database implementations
│   └── usecase/        # Business logic implementations
└── test/               # Integration and E2E tests
```

## 📄 License

This project is licensed under the MIT License.

