# Document Reader Chatbot - Backend

A scalable Go backend service for a document reader chatbot application built with Clean Architecture principles, featuring user management, document storage, and AI-powered chat functionality.

## 🏗️ Architecture

This backend follows **Clean Architecture** principles with clear separation of concerns:

### Project Structure

```
backend/
├── cmd/
│   └── server/           # Application entry point
├── configs/              # Configuration management
├── internal/             # Core application logic
│   ├── user/            # User domain service
│   │   ├── models.go    # Domain models
│   │   ├── repository.go # Data access layer
│   │   ├── service.go   # Business logic
│   │   ├── controller.go # HTTP handlers
│   │   └── router.go    # Route definitions
│   └── document/        # Document domain service
│       ├── models.go    # Domain models
│       ├── repository.go # Data access layer
│       ├── service.go   # Business logic
│       ├── controller.go # HTTP handlers
│       └── router.go    # Route definitions
├── pkg/                 # Shared utilities
│   ├── database/        # Database connection
│   ├── middleware/      # HTTP middleware
│   ├── observability/   # OpenTelemetry setup
│   └── utils/           # Utility functions
└── test/               # Test utilities and mocks
```

### Key Features

- **Clean Architecture** with clear layer separation
- **Domain-Driven Design** with feature-based organization
- **OpenTelemetry** integration for observability
- **JWT Authentication** with role-based access control
- **MongoDB** with optimized queries and indexing
- **Comprehensive Error Handling** with structured responses
- **Rate Limiting** and security middleware
- **Graceful Shutdown** with proper resource cleanup

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgresqlPostgresql
- Docker (optional, for containerized dependencies)

### Setup

1. **Clone and navigate to the backend directory**
   ```bash
   git clone <repository-url>
   cd document-reader-chatbot/backend
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Install development tools**
   ```bash
   make setup-tools
   ```

4. **Setup environment**
   ```bash
   make setup
   # Edit .env file with your configuration
   ```

5. **Start MongoDB**
   ```bash
   make db-start
   ```

6. **Run the application**
   ```bash
   make run
   # or for development with hot reload
   make dev
   ```

The server will start on `http://localhost:8080`

## 🔧 Configuration

Copy `config.env.example` to `.env` and configure the following:

### Server Configuration
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/production)
- `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT` - Server timeouts

### Database Configuration
- `MONGODB_URI` - MongoDB connection string
- `DATABASE_NAME` - Database name
- `DB_MAX_POOL_SIZE`, `DB_MIN_POOL_SIZE` - Connection pool settings

### JWT Configuration
- `JWT_SECRET_KEY` - Secret key for JWT signing
- `JWT_EXPIRATION_HOURS` - Token expiration time
- `JWT_ISSUER` - Token issuer

### Observability Configuration
- `JAEGER_ENDPOINT` - Jaeger tracing endpoint
- `METRICS_PORT` - Prometheus metrics port

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### User Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Get Profile
```http
GET /profile
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

### Document Endpoints

#### Create Document
```http
POST /documents
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Document",
  "description": "Document description",
  "content": "Document content...",
  "file_name": "document.txt",
  "file_type": "text/plain",
  "file_size": 1024,
  "is_public": false,
  "tags": ["tag1", "tag2"]
}
```

#### Get Document
```http
GET /documents/{id}
Authorization: Bearer <token>
```

#### List Public Documents
```http
GET /documents?page=1&limit=10&search=query
```

#### List My Documents
```http
GET /documents/my?page=1&limit=10
Authorization: Bearer <token>
```

#### Update Document
```http
PUT /documents/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "is_public": true
}
```

#### Delete Document
```http
DELETE /documents/{id}
Authorization: Bearer <token>
```

### Chat Endpoints

#### Create Chat Session
```http
POST /documents/{id}/chat
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Chat about document"
}
```

#### List Chat Sessions
```http
GET /chat?page=1&limit=10
Authorization: Bearer <token>
```

#### Get Chat Session
```http
GET /chat/{sessionId}
Authorization: Bearer <token>
```

#### Send Message
```http
POST /chat/{sessionId}/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "What is this document about?"
}
```

#### Delete Chat Session
```http
DELETE /chat/{sessionId}
Authorization: Bearer <token>
```

## 🧪 Testing

### Run Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Run Tests with Race Detection
```bash
make test-race
```

### Run Benchmarks
```bash
make benchmark
```

## 🔍 Code Quality

### Lint Code
```bash
make lint
```

### Format Code
```bash
make fmt
```

### Security Check
```bash
make security
```

### Vet Code
```bash
make vet
```

## 📊 Observability

### Monitoring

The application includes comprehensive observability:

- **Distributed Tracing** with OpenTelemetry and Jaeger
- **Metrics** with Prometheus
- **Structured Logging** with trace correlation
- **Health Checks** at `/health`

### Start Jaeger
```bash
make jaeger-start
```

Access Jaeger UI at `http://localhost:16686`

### Metrics

Prometheus metrics are exposed at `http://localhost:9090/metrics`

## 🐳 Docker

### Build Docker Image
```bash
make build-docker
```

### Start All Services
```bash
make docker-up
```

### Stop All Services
```bash
make docker-down
```

## 🏗️ Development

### Project Guidelines

1. **Follow Clean Architecture** - Keep layers separated
2. **Use Interfaces** - All public functions should work with interfaces
3. **Handle Errors Explicitly** - Use wrapped errors for traceability
4. **Add Tracing** - Include OpenTelemetry spans for operations
5. **Write Tests** - Maintain good test coverage
6. **Document Code** - Use GoDoc style comments

### Adding New Features

1. **Define Domain Models** in `models.go`
2. **Create Repository Interface** in `repository.go`
3. **Implement Business Logic** in `service.go`
4. **Add HTTP Handlers** in `controller.go`
5. **Setup Routes** in `router.go`
6. **Add Tests** for all layers
7. **Update Documentation**

### Database Patterns

- Use MongoDB aggregation pipelines for complex queries
- Implement proper indexing for performance
- Use transactions for operations requiring atomicity
- Follow consistent naming conventions

## 🚀 Deployment

### Build for Production
```bash
make deploy-build
```

### Environment Variables

Ensure all required environment variables are set in production:
- Use strong JWT secret keys
- Set proper database connection strings
- Configure observability endpoints
- Set security headers appropriately

### Security Considerations

- Use HTTPS in production
- Implement rate limiting
- Validate all inputs
- Use secure password hashing
- Implement proper CORS policies
- Keep dependencies updated

## 🤝 Contributing

1. Follow the existing code structure
2. Write tests for new functionality
3. Run linting and formatting before committing
4. Update documentation for new features
5. Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

For questions or issues:

1. Check the documentation
2. Review existing issues
3. Create a new issue with detailed information
4. Include logs and error messages

## 🗂️ Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [MongoDB Go Driver](https://docs.mongodb.com/drivers/go/)
- [Gin Web Framework](https://gin-gonic.com/docs/) 