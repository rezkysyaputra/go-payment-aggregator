# Go Payment Aggregator

A robust, cloud-ready Payment Aggregator service dealing with multiple payment gateways (Midtrans, Xendit) through a unified API interface. Built with Golang, Clean Architecture, and reliability in mind.

## 🚀 Features

*   **Multi-Gateway Support**: Seamlessly switch between **Midtrans**, **Xendit**, and **Mock** (Local Testing) providers.
*   **Unified API**: Singe `CreateTransaction` endpoint handles logic for different providers intelligently.
*   **Generic Webhook Handling**: Standardized webhook processing pipeline for all providers.
*   **Merchant Notifications**: Automatically sends webhook notifications back to the merchant's `callback_url` upon payment status updates.
*   **Duplicate Order Protection**: Prevents double-spending or duplicate transaction creation for the same Order ID.
*   **Dockerized**: Ready to deploy with `docker-compose` including PostgreSQL database.
*   **Modular Architecture**: Built using Hexagonal / Clean Architecture principles for maintainability.

## 🛠 Tech Stack

*   **Language**: Go (Golang) 1.22+
*   **Framework**: Gin Gonic
*   **Database**: PostgreSQL
*   **ORM**: GORM
*   **Configuration**: Viper
*   **Logging**: Logrus
*   **Testing**: Go Test + Testify + Mocking

## 📡 API Documentation

The full OpenAPI (Swagger) specification is available in `api/openapi.json`.

### Core Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/v1/merchant/register` | Register a new merchant & get API Key |
| `POST` | `/v1/transaction` | Create a new payment transaction |
| `GET` | `/v1/transaction/{id}` | Get transaction details by System ID |
| `GET` | `/v1/transaction/order/{order_id}` | Get transaction details by Order ID |

### Supported Providers
*   `midtrans`: Uses Midtrans Snap
*   `xendit`: Uses Xendit Invoice
*   `mock`: Local simulation for testing

## 💻 Getting Started

### Prerequisites
*   Docker & Docker Compose
*   Go 1.22+ (for local development)

### Configuration
1.  Copy `config_example.json` to `config.docker.json` (for Docker) or `config.json` (for Local).
    ```bash
    cp config_example.json config.json
    ```
2.  Fill in your API Keys:
    ```json
    "midtrans": { "server_key": "YOUR_SB_KEY" },
    "xendit": { "api_key": "YOUR_SECRET_KEY", "callback_token": "YOUR_TOKEN" }
    ```

### Running with Docker (Recommended)
```bash
docker-compose up -d --build
```
The API will be available at `http://localhost:8080`.

### Running Locally
1.  Start Database:
    ```bash
    docker-compose up -d db migrate
    ```
2.  Run App:
    ```bash
    go run cmd/server/main.go
    ```

## 🧪 Testing
Run unit tests including provider validations:
```bash
go test ./test/... -v
```

## 📂 Project Structure
```
.
├── api/             # API Specifications (OpenAPI)
├── cmd/             # Entrypoints
├── internal/
│   ├── config/      # Configuration & Bootstrap
│   ├── domain/      # Business Logic (Entities, Repos, Services)
│   ├── gateway/     # External Payment Gateway Implementations (Strategy Pattern)
│   ├── handler/     # HTTP Handlers
│   ├── middleware/  # Auth Middleware
│   └── router/      # Gin Route Definitions
├── migrations/      # SQL Migrations
└── test/            # Integration/Unit Tests
```
