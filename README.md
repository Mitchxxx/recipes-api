# Recipes API

The Recipes API is a RESTful service for managing recipes. It allows users to perform common operations like creating, reading, updating, and deleting recipes. The API provides multiple authentication implementations including JWT-based authentication, Auth0 integration, and session-based authentication.

## Features

* User authentication and authorization with multiple implementation options
* CRUD operations for managing recipes
* Redis caching for improved performance
* MongoDB for data persistence
* Session storage using Redis (optional)
* Designed with scalability and security best practices
* Well-documented endpoints using Swagger/OpenAPI
* Graceful shutdown handling

## Technology Stack

* **Language**: Go 1.24.5
* **Web Framework**: Gin
* **Database**: MongoDB
* **Caching**: Redis
* **Authentication**: JWT / Auth0 / Sessions
* **Password Hashing**: bcrypt

## Prerequisites

* Go 1.24.5 or higher
* MongoDB instance
* Redis instance
* Auth0 account (for auth0 implementation only)

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/Mitchxxx/recipes-api.git
cd recipes-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following variables:
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=recipes_db
REDIS_ADDRESS=localhost:6379
JWT_SECRET=your-secret-key
X_API_KEY=your-api-key
SESSION_SECRET=your-session-secret
REDIS_PASSWORD=
# For Auth0 implementation only:
AUTH_DOMAIN=your-auth0-domain.auth0.com
AUTH0_API_IDENTIFIER=your-api-identifier
```

4. Run the application:
```bash
# Main JWT implementation
go run main.go

# Or use Air for hot reload
air

# For Auth0 implementation
cd auth0 && go run main.go

# For Cookies/Session implementation
cd cookies && go run main.go

# For API implementation (JWT with sessions)
cd api && go run main.go
```

The server will start on `http://localhost:8080`

## Authentication Implementations

This repository includes multiple authentication implementations:

1. **Main (JWT)** - Located in root directory, uses JWT tokens with username/password authentication
2. **Auth0** - Located in `/auth0`, integrates with Auth0 for authentication
3. **Cookies** - Located in `/cookies`, uses session-based authentication with Redis
4. **API** - Located in `/api`, combines JWT with session storage

### Authentication Flow

Users must sign in to obtain authentication credentials and include them in requests for protected endpoints (e.g., creating, updating, or deleting recipes).

# Endpoints


## Authentication

### POST /signin
Sign in a user with their credentials and generate a JWT token.

**Request Body:**
```json
{
    "username": "admin",
    "password": "password123"
}
```

**Response (200 OK):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "2025-08-01T16:45:30Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid username or password

### POST /refresh
Refresh an existing JWT token.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "2025-08-01T18:00:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or expired token
- `400 Bad Request`: Token is not expired yet

## Recipes

### GET /recipes
Retrieve all recipes (public endpoint - no authentication required).

**Response (200 OK):**
```json
[
    {
        "id": "12345abcd",
        "name": "Tomato Soup",
        "tags": ["soup", "vegetarian"],
        "ingredients": ["Tomatoes", "Onion", "Garlic"],
        "instructions": ["Chop vegetables", "Boil water", "Add ingredients"],
        "publishedAt": "2023-10-05T15:00:00Z"
    }
]
```

### POST /recipes
Create a new recipe (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
    "name": "Tomato Soup",
    "tags": ["soup", "vegetarian"],
    "ingredients": ["Tomatoes", "Onion", "Garlic"],
    "instructions": ["Chop vegetables", "Boil water", "Add ingredients"]
}
```

**Response (200 OK):**
```json
{
    "message": "Recipe created successfully",
    "recipe": {
        "id": "12345abcd",
        "name": "Tomato Soup",
        "tags": ["soup", "vegetarian"],
        "ingredients": ["Tomatoes", "Onion", "Garlic"],
        "instructions": ["Chop vegetables", "Boil water", "Add ingredients"],
        "publishedAt": "2023-10-05T15:00:00Z"
    }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid authentication token


### GET /recipes/{id}
Retrieve a specific recipe by its ID (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (string): The ID of the recipe

**Response (200 OK):**
```json
{
    "id": "12345abcd",
    "name": "Tomato Soup",
    "tags": ["soup", "vegetarian"],
    "ingredients": ["Tomatoes", "Onion", "Garlic"],
    "instructions": ["Chop vegetables", "Boil water", "Add ingredients"],
    "publishedAt": "2023-10-05T15:00:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid authentication token
- `404 Not Found`: Recipe not found

### PUT /recipes/{id}
Update an existing recipe by its ID (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (string): The ID of the recipe

**Request Body:**
```json
{
    "name": "Updated Tomato Soup",
    "tags": ["soup", "healthy"],
    "ingredients": ["Tomatoes", "Spinach"],
    "instructions": ["Blend vegetables", "Heat", "Serve"]
}
```

**Response (200 OK):**
```json
{
    "message": "recipe has been updated"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid authentication token
- `404 Not Found`: Recipe not found

### DELETE /recipes/{id}
Delete a specific recipe by its ID (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (string): The ID of the recipe

**Response (200 OK):**
```json
{
    "message": "Recipe has been deleted"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid authentication token
- `404 Not Found`: Recipe not found
- `500 Internal Server Error`: Server-side issues

## Error Responses

Common error responses across endpoints:

* **400 Bad Request**: Invalid input or malformed request
* **401 Unauthorized**: Authentication required or token is invalid/expired
* **404 Not Found**: Resource not found (e.g., recipe doesn't exist)
* **500 Internal Server Error**: Server-side issues

## Performance

The API implements Redis caching for the GET /recipes endpoint to improve performance. Benchmark results show significant performance improvements when caching is enabled.

## Development

### Project Structure

```
.
├── main.go                 # Main JWT implementation
├── handlers/              # Request handlers
├── models/                # Data models
├── api/                   # JWT with sessions implementation
├── auth0/                 # Auth0 integration
├── cookies/               # Session-based authentication
├── recipes.json           # Sample recipe data
└── swagger.json           # API documentation
```

### Running with Hot Reload

The project includes Air configuration for hot reloading during development:

```bash
air
```

### API Documentation

Swagger/OpenAPI documentation is available in `swagger.json`. You can view it using any Swagger UI tool.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is part of the Building Distributed Applications in Gin course.

## Contact

Mitchel Egboko - megboko@gmail.com

Project Link: [https://github.com/Mitchxxx/recipes-api](https://github.com/Mitchxxx/recipes-api)