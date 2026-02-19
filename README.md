# Go Gin GORM Boilerplate 🚀

## 📝 Description

**Go Gin GORM Boilerplate** is a complete template for building robust, scalable, and production-ready REST API backends using Go. This boilerplate implements **Clean Architecture** with modern best practices in Go backend development.

The project comes equipped with various enterprise-level features such as authentication, payment gateway integration (Tripay), email verification, AWS S3 integration, and much more.

---

## ✨ Main Feature

### 🔐 Authentication & Authorization
- **JWT-based Authentication**: Secure token-based authentication using JWT
- **User Registration & Login**: Complete registration system with validation
- **Email Verification**: Automated email verification with AES encrypted token
- **Password Management**: Secure forgot password and reset password functionality
- **Role-based Access Control**: Support for multi-role (Admin, User)

### 💳 Payment Integration
- **Tripay Payment Gateway**: Complete integration with Tripay for multiple payment methods
- **HMAC-SHA256 Signature Verification**: Secure webhook with signature verification
- **Transaction Management**: Tracking transaction status (PAID, FAILED, EXPIRED, REFUND)
- **Invoice Generation**: Generate invoice URL for payment

### 📧 Email Services
- **Email Verification**: Automated email verification with HTML templates
- **Forgot Password**: Send secure password reset emails with tokens
- **SMTP Integration**: Support SMTP with Gmail and other email providers

### ☁️ Cloud Storage
- **AWS S3 Integration**: Upload and manage files to AWS S3
- **Secure File Upload**: File upload with encryption and validation

### 🛠 Advanced Features
- **Clean Architecture**: Separation of concerns with layers: Entity, DTO, Repository, Service, Controller
- **Database Migration**: Automatic migration with GORM
- **Seeder**: Seed data for development
- **CORS Middleware**: Pre-configured CORS for frontend integration
- **Error Handling**: Centralized error handling with custom error messages
- **Logging**: Comprehensive logging with Logrus
- **Data Validation**: Input validation in DTO layer
- **Pagination Support**: Built-in pagination utility

---

## � Project Structure

```
go-gin-gorm-boilerplate/
├── cmd/                          # CLI Commands
│   └── command.go               # Command handler (migrate, seed)
│
├── constants/                    # Application constants
│   └── common.go
│
├── database/                     # Database setup & migrations
│   ├── database.go              # Connection setup
│   ├── migrations/              # Database migrations
│   │   └── migrate.go
│   └── seeders/                 # Data seeders
│       ├── seeder.go
│       ├── data/
│       ├── json/
│       └── seeds/
│           └── user_seed.go
│
├── internal/                     # Internal packages (not exported)
│   ├── api/                      # API layer
│   │   ├── controller/           # HTTP Handlers
│   │   │   ├── transaction_controller.go
│   │   │   └── user_controller.go
│   │   ├── repository/           # Data Access Layer
│   │   │   ├── common.go
│   │   │   ├── transaction_repository.go
│   │   │   └── user_repository.go
│   │   ├── routes/               # Route definitions
│   │   │   ├── transaction_route.go
│   │   │   └── user_route.go
│   │   └── service/              # Business Logic Layer
│   │       ├── jwt_service.go
│   │       ├── tansaction_service.go
│   │       └── user_service.go
│   │
│   ├── config/                   # Configuration & Setup
│   │   ├── rest_config.go        # REST API initialization & DI
│   │   └── rest_router_config.go # Router setup & middleware
│   │
│   ├── dto/                      # Data Transfer Objects
│   │   ├── common.go
│   │   ├── transaction_dto.go
│   │   └── user_dto.go
│   │
│   ├── entity/                   # Domain Models
│   │   ├── common.go
│   │   ├── transaction_entity.go
│   │   └── user_entity.go
│   │
│   ├── middleware/               # HTTP Middleware
│   │   ├── authentication.go
│   │   ├── cors.go
│   │   ├── only_allow.go
│   │   └── time.go
│   │
│   ├── pkg/                      # Reusable packages
│   │   ├── logger/               # Logging
│   │   │   └── logger.go
│   │   ├── mailer/               # Email service
│   │   │   ├── mailer.go
│   │   │   ├── makeMail.go
│   │   │   └── template/
│   │   │       ├── forgot_password_email.html
│   │   │       └── verification_email.html
│   │   ├── pagination/           # Pagination utilities
│   │   │   ├── conv.go
│   │   │   └── meta.go
│   │   ├── payment/              # Payment integrations
│   │   │   └── tripay/
│   │   │       ├── client.go
│   │   │       ├── signature.go
│   │   │       └── tripay.go
│   │   ├── response/             # Response formatting
│   │   │   └── response.go
│   │   └── storage/              # Cloud storage
│   │       └── aws_s3.go
│   │
│   └── utils/                    # Utility functions
│       ├── aes.go                # Encryption utilities
│       ├── file.go
│       ├── is_prod.go
│       └── password.go
│
├── assets/                      # Static files & images
├── .env.example                 # Environment variables template
├── Dockerfile                   # Docker image configuration
├── docker-compose.yml          # Docker Compose setup
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── main.go                      # Application entry point
└── README.md                    # This file
```

### Folder Descriptions

| Folder | Purpose |
|--------|---------|
| **cmd/** | CLI commands for database operations (migrate, seed) |
| **constants/** | Application-wide constants |
| **database/** | Database connection, migrations, and seeders |
| **internal/api/controller/** | HTTP request handlers following Clean Architecture |
| **internal/api/service/** | Business logic layer, orchestrates operations between controllers and repositories |
| **internal/api/repository/** | Data access layer, abstracts database operations |
| **internal/api/routes/** | Route definitions and endpoint configurations |
| **internal/config/** | Core configuration: REST API setup, dependency injection, router initialization |
| **internal/dto/** | Data Transfer Objects for request/response serialization |
| **internal/entity/** | Domain models representing database tables |
| **internal/middleware/** | HTTP middleware for authentication, CORS, logging, etc. |
| **internal/pkg/logger/** | Logging utility |
| **internal/pkg/mailer/** | Email service with SMTP integration |
| **internal/pkg/pagination/** | Pagination utilities |
| **internal/pkg/payment/tripay/** | Tripay payment gateway integration |
| **internal/pkg/response/** | Standardized API response formatting |
| **internal/pkg/storage/aws_s3/** | AWS S3 cloud storage integration |
| **internal/utils/** | Utility functions: encryption (AES), file handling, environment checks, password hashing |
| **routes/** | Route definitions and endpoint configurations |
| **constants/** | Application-wide constants |
| **helpers/** | Helper functions for common operations |

### Configuration Flow

```
main.go
  ↓
database.SetUpDatabaseConnection() 
  ↓
config.NewRestConfig(db)  [internal/config/rest_config.go]
  ├─ Dependency Injection
  ├─ Initialize Services, Repositories, Controllers
  └─ NewRouter(app)      [internal/config/rest_router_config.go]
      ├─ Setup Middleware
      ├─ Register Routes
      └─ Configure Static Files
  ↓
restConfig.Start()
  ↓
Graceful Shutdown (30s timeout)
```

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|-----------|
| **Framework** | Gin Web Framework |
| **Database** | PostgreSQL with GORM ORM |
| **Authentication** | JWT (golang-jwt) |
| **Encryption** | AES + bcrypt |
| **Email** | Gomail (SMTP) |
| **Cloud Storage** | AWS SDK v2 (S3) |
| **Payment** | Tripay API |
| **Notifications** | Discord Webhook |
| **Logging** | Logrus |
| **Deployment** | Docker & Docker Compose |
| **Language** | Go 1.24.4 |

---

## 🔧 Configuration Architecture

### REST API Configuration (`internal/config/`)

The configuration package handles all REST API setup with a clean, modular approach:

#### rest_config.go
- **Purpose**: Initialize REST API and manage dependency injection
- **Responsibility**:
  - Create database connection
  - Instantiate all services (JWT, Mailer)
  - Create all repositories
  - Create all controllers
  - Register all routes
  - Start HTTP server with graceful shutdown support
- **Key Methods**:
  - `NewRestConfig(db)` - Initialize REST configuration
  - `Start()` - Start the HTTP server
  - `Shutdown(ctx)` - Graceful shutdown with timeout

#### rest_router_config.go
- **Purpose**: Configure router, middleware, and endpoints
- **Responsibility**:
  - Setup Gin router
  - Apply global middleware (CORS, logging)
  - Configure no-route handlers
  - Setup health check endpoints
  - Register static file servers
- **Key Functions**:
  - `NewRouter(server)` - Configure and return router

### Dependency Injection Flow

```
NewRestConfig(db)
├─ JWT Service
├─ Mailer Service
├─ Repositories (User, Transaction)
├─ Services (User, Transaction)
├─ Controllers (User, Transaction)
└─ Routes Registration
```

---

### Clean Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│              Application Entry Point                │
│ main.go → Database Setup → config.NewRestConfig()  │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│  Configuration Layer (internal/config/)             │
│  - REST API initialization                          │
│  - Dependency Injection (DI)                        │
│  - Router & Middleware Setup                        │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│      HTTP Layer (internal/api/controller/)          │
│  - Handle HTTP requests/responses                   │
│  - Input validation via DTO                         │
│  - Response formatting                              │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│    Business Logic Layer (internal/api/service/)     │
│  - Orchestrate core business rules                  │
│  - Data transformation & validation                 │
│  - Multi-step workflows                             │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│   Data Access Layer (internal/api/repository/)      │
│  - Abstract database operations                     │
│  - GORM query builder                               │
│  - Data persistence & retrieval                     │
└────────────────┬────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────┐
│        Database Layer (PostgreSQL/GORM)             │
│  - Data storage                                     │
│  - Transaction management                           │
└─────────────────────────────────────────────────────┘

Cross-Cutting Concerns:
├─ internal/middleware/ (Authentication, CORS, Logging)
├─ internal/pkg/ (Email, Payment Gateway, Logging, Response Formatting)
├─ internal/utils/ (Encryption, Password Hashing, File Handling, Environment checks)
└─ constants/ (App-wide constants)
```

---

## 📦 Installation

### Prerequisites
Make sure you have installed:
- **Go 1.24.4 or newer** - [Download](https://golang.org/dl/)
- **PostgreSQL 12 or newer** - [Download](https://www.postgresql.org/download/)
- **Git** - [Download](https://git-scm.com/)
- **Docker & Docker Compose** (for development with containers) - [Download](https://docs.docker.com/get-docker/)

### Step 1: Clone Repository
```bash
git clone https://github.com/Shabrinashsf/go-gin-gorm-boilerplate.git
cd go-gin-gorm-boilerplate
```

### Step 2: Setup Environment Variables
```bash
# Copy .env.example to .env
cp .env.example .env

# Edit .env with appropriate values
nano .env  # or use your favorite editor
```

### Step 3: Install Dependencies
```bash
# Download all Go dependencies
go mod download

# or
go get ./...
```

### Step 4: Setup Database

#### Option A: Local PostgreSQL
```bash
# Run migration to create tables
go run main.go --migrate

# (Optional) Seed data to the database
go run main.go --seed
```

#### Option B: Docker PostgreSQL
```bash
# Run PostgreSQL container
docker-compose up -d postgres

# Then run migration
go run main.go --migrate
go run main.go --seed
```

### Step 5: Run Application
```bash
# Development mode
go run main.go

# Production mode (build first)
go build -o main .
./main
```

The application will run at `http://localhost:8888`

---

## ⚙️ Configuration

### Setup SMTP Gmail

1. **Open Google Account Security**
   - Login to [myaccount.google.com](https://myaccount.google.com)
   - Select "Security" from the left sidebar
   - Enable "2-Step Verification" if not already enabled

2. **Generate App Password**
   - Click "App passwords" (only appears if 2FA is enabled)
   - Select "Mail" and "Windows Computer"
   - Google will generate a 16-character password
   - Copy this password to `SMTP_AUTH_PASSWORD` in .env


### Setup AWS S3

1. **Create IAM User**
   - Login ke [AWS Console](https://console.aws.amazon.com/)
   - Open IAM → Users → Create User
   - Enable "Programmatic access"
   - Attach policy: `AmazonS3FullAccess`
   - Save Access Key ID and Secret Access Key

2. **Create S3 Bucket**
   - Open S3 Console
   - Create bucket with a unique name
   - Configure bucket policies for public read (if needed)

### Setup Tripay Payment Gateway

1. **Create Tripay Account**
   - Register at [Tripay Dashboard](https://dashboard.tripay.co.id/)
   - Verify your email and login

2. **Get API Credentials**
   - Login to Tripay dashboard
   - Open Settings → API Keys
   - Copy: Merchant Code, API Key, and Private Key

3. **Setup Webhook**
   - Open Settings → Webhook Configuration
   - Callback URL: `https://yourapp.com/api/transaction/webhook/tripay`
   - Event: Payment Status
   - Save configuration

---

## 🚀 Usage Guide

### 1. Running Local Development

```bash
# Terminal 1: Start PostgreSQL
docker-compose up -d postgres

# Terminal 2: Run application
go run main.go
```

Application is ready at `http://localhost:8888`

### 2. Graceful Shutdown

The application implements graceful shutdown with a 30-second timeout:

```bash
# Start application
go run main.go

# To shutdown gracefully (Ctrl+C or SIGTERM)
# 1. Stops accepting new requests
# 2. Waits for in-flight requests to complete (max 30 seconds)
# 3. Closes HTTP server
# 4. Closes database connections
# 5. Exits cleanly
```

On interrupt signal, you'll see:
```
[info] Received signal: interrupt
[info] Starting graceful shutdown...
[info] Shutting down HTTP server...
[info] HTTP server stopped
[info] Graceful shutdown completed
[info] Application exited
```

### 3. Database Management

```bash
# Run migration to create tables
go run main.go --migrate

# Seed example data
go run main.go --seed

# View help commands
go run main.go --help
```

### 4. API Testing with Bruno

Complete API documentation is available at:
- 📍 **Bruno Documentation**: [docs-go-gin-gorm-boilerplate](https://github.com/Shabrinashsf/docs-go-gin-gorm-boilerplate)

Download Bruno API Client: https://www.usebruno.com/

```bash
# Clone documentation repository
git clone https://github.com/Shabrinashsf/docs-go-gin-gorm-boilerplate.git

# Open in Bruno API client
# 1. Download Bruno: https://www.usebruno.com/
# 2. File → Open Collection
# 3. Select folder docs-go-gin-gorm-boilerplate
# 4. Set environment in Bruno
# 5. Test all endpoints
```

---

## 🐳 Docker 

### Local Development with Docker

```bash
# Build and run all services
docker-compose up --build

# Run migration and seed in container
docker-compose exec app go run main.go --migrate
docker-compose exec app go run main.go --seed

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app
docker-compose logs -f postgres
```

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Tripay](https://tripay.co.id/)
- [AWS SDK Go](https://github.com/aws/aws-sdk-go-v2)

---

<div align="center">

**[⬆ Back to Top](#go-gin-gorm-boilerplate-)**

Made with ❤️ by [Shabrinashsf](https://github.com/Shabrinashsf)

</div>