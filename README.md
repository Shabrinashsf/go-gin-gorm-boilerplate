# Go Gin GORM Boilerplate 🚀

<p align="center">
  <a href="#main-feature">Feature</a> •
  <a href="#project-structure">Project Structure</a> •
  <a href="#installation">Installation</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#api-documentation">API</a>
</p>

---

## 📋 Daftar Isi
- [Description](#description)
- [Main Feature](#main-feature)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage Guide](#usage-guide)
- [Docker](#docker--deployment)

---

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

### Clean Architecture Diagram

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Controller)         │
│  - Handle HTTP requests/responses       │
│  - Input validation                     │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│         Business Logic (Service)        │
│  - Core business rules                  │
│  - Data transformation                  │
│  - Workflow orchestration               │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│        Data Access (Repository)         │
│  - Database operations                  │
│  - Data persistence                     │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│      Database (PostgreSQL/GORM)        │
│  - Data storage                         │
└─────────────────────────────────────────┘
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

### 2. Database Management

```bash
# Run migration to create tables
go run main.go --migrate

# Seed example data
go run main.go --seed

# View help commands
go run main.go --help
```

### 3. API Testing with Bruno

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