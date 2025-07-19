# Authentication API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication Endpoints

### 1. Register User

**Endpoint:** `POST /auth/register`

**Description:** Create a new user account

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john.doe@example.com", 
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Request Validation:**
- `username`: Required, 3-50 characters, unique
- `email`: Required, valid email format, unique
- `password`: Required, minimum 8 characters
- `first_name`: Required, 2-50 characters
- `last_name`: Required, 2-50 characters

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-08T10:00:00Z",
    "updated_at": "2024-01-08T10:00:00Z"
  }
}
```

**Error Responses:**

**400 Bad Request - Invalid Data:**
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request data",
    "details": "Username is required"
  }
}
```

**409 Conflict - User Exists:**
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT", 
    "message": "user with email john.doe@example.com already exists"
  }
}
```

### 2. User Login

**Endpoint:** `POST /auth/login`

**Description:** Authenticate user and get JWT token

**Request Body:**
```json
{
  "email_or_username": "john.doe@example.com",
  "password": "password123"
}
```

**Note:** You can use either email or username in the `email_or_username` field:
- If the field contains `@`, it will be treated as an email
- Otherwise, it will be treated as a username

**Request Validation:**
- `email_or_username`: Required
- `password`: Required

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe", 
      "email": "john.doe@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-08T10:00:00Z",
      "updated_at": "2024-01-08T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**

**400 Bad Request - Invalid Data:**
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request data",
    "details": "Password is required"
  }
}
```

**401 Unauthorized - Invalid Credentials:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "invalid credentials"
  }
}
```

**401 Unauthorized - Account Disabled:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "user account is disabled"
  }
}
```

## Using JWT Token

After successful login, use the returned JWT token in the Authorization header for protected endpoints:

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/profile
```

## Token Information

- **Expiration:** 24 hours (configurable via `JWT_EXPIRATION_HOURS`)
- **Algorithm:** HS256
- **Claims Include:** user_id, email, role

## Database Schema

The authentication system uses the following PostgreSQL tables:

### users table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
```

## Security Features

- ✅ **Password Hashing:** bcrypt with default cost
- ✅ **JWT Tokens:** Signed with HS256 
- ✅ **Input Validation:** Comprehensive request validation
- ✅ **Rate Limiting:** Built-in middleware protection
- ✅ **CORS:** Configurable cross-origin support
- ✅ **Security Headers:** Standard security headers
- ✅ **Unique Constraints:** Username and email uniqueness
- ✅ **Soft Deletes:** User accounts can be soft deleted
- ✅ **Role-based Access:** User/Admin roles supported 