# OTP Authentication Service

A comprehensive backend service in Golang that implements OTP-based login and registration, along with basic user management features.

## Features

- **OTP-based Authentication**: Secure one-time password authentication
- **Rate Limiting**: Prevents abuse with configurable limits
- **JWT Tokens**: Secure session management
- **User Management**: CRUD operations with pagination and search
- **Swagger Documentation**: Complete API documentation
- **Docker Support**: Easy deployment with Docker and docker-compose
- **Redis Integration**: Fast and reliable OTP storage

## Database Choice: Redis

**Why Redis?**

1. **Performance**: Redis is an in-memory data store, providing extremely fast read/write operations
2. **TTL Support**: Built-in expiration mechanism perfect for OTP storage (2-minute expiry)
3. **Rate Limiting**: Atomic operations for implementing rate limiting efficiently
4. **Simplicity**: No complex schema management, perfect for temporary data
5. **Scalability**: Can handle high concurrent OTP requests
6. **Persistence**: Optional persistence for data recovery

## Quick Start

### Option 1: Local Development

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd otp-auth-service
   go mod download
   ```

2. **Start Redis** (using Docker):
   ```bash
   docker run -d -p 6379:6379 redis:7-alpine
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Access the service**:
   - API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html

### Option 2: Docker Compose (Recommended)

1. **Start all services**:
   ```bash
   docker-compose up -d
   ```

2. **View logs**:
   ```bash
   docker-compose logs -f app
   ```

## API Examples

### Request OTP
```bash
curl -X POST http://localhost:8080/api/v1/auth/request-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'
```

### Verify OTP
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "otp": "123456"}'
```

### Get Users (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Application port |
| `REDIS_ADDR` | `localhost:6379` | Redis server address |
| `REDIS_PASSWORD` | `` | Redis password (if any) |
| `REDIS_DB` | `0` | Redis database number |
| `JWT_SECRET` | `your-secret-key-change-in-production` | JWT signing secret |

## Security Features

1. **OTP Expiration**: OTPs expire after 2 minutes
2. **Rate Limiting**: Maximum 3 OTP requests per phone number per 10 minutes
3. **JWT Tokens**: 24-hour expiry with secure signing
4. **Input Validation**: Comprehensive request validation
5. **CORS Protection**: Configurable cross-origin resource sharing
