# Recipes API
The Recipes API is a RESTful service for managing recipes. It allows users to perform common operations like creating, reading, updating, and deleting recipes. The API is secured using JWT tokens generated via an Auth0 integration.

# Features
* User authentication and authorization using JWT (via Auth0).
* CRUD operations for managing recipes.
* Session storage using Redis for efficient state management (optional).
* Designed with scalability and security best practices.
* Well-documented endpoints using Swagger/OpenAPI.

## Authentication

Users must sign in to obtain a JWT (/signin) and include it in the Authorization header (Bearer <token>) for protected endpoints (e.g., creating or updating recipes).

# Endpoints


## Authentication

### POST /signin
Sign in a user with their credentials and generate a JWT.

Request payload

```
{
    "username": "admin",
    "password": "password123"
}
```
 Response

```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "2025-08-01T16:45:30Z"
}

```

### POST /Refresh

* Headers

```
  Authorization: Bearer <token>
```

 Response

```
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "2025-08-01T18:00:00Z"
  }


```

## Recipes

### GET /recipes
Retrieve all recipes

 Response

```
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

### POST /signin
Create a new recipe. (Requires authentication)

Request payload

```
  {
    "name": "Tomato Soup",
    "tags": ["soup", "vegetarian"],
    "ingredients": ["Tomatoes", "Onion", "Garlic"],
    "instructions": ["Chop vegetables", "Boil water", "Add ingredients"]
  }

```
 Response

```
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


### GET /recipes/{id}
Retrieve a specific recipe by its id.

* Path Parameter:

```
id: The ID of the recipe.

```

 Response

```
{
    "id": "12345abcd",
    "name": "Tomato Soup",
    "tags": ["soup", "vegetarian"],
    "ingredients": ["Tomatoes", "Onion", "Garlic"],
    "instructions": ["Chop vegetables", "Boil water", "Add ingredients"],
    "publishedAt": "2023-10-05T15:00:00Z"
}

```

### PUT /recipes/{id}
Update an existing recipe by its id. (Requires authentication)

* Path Parameter:

```
id: The ID of the recipe.

```

Request payload

```
{
    "name": "Updated Tomato Soup",
    "tags": ["soup", "healthy"],
    "ingredients": ["Tomatoes", "Spinach"],
    "instructions": ["Blend vegetables", "Heat", "Serve"]
}


```
 Response

```
{
    "message": "Recipe updated successfully",
    "recipe": {
      "id": "12345abcd",
      "name": "Updated Tomato Soup",
      "tags": ["soup", "healthy"],
      "ingredients": ["Tomatoes", "Spinach"],
      "instructions": ["Blend vegetables", "Heat", "Serve"],
      "publishedAt": "2023-10-05T15:45:00Z"
    }
}

```

### DELETE /recipes/{id}
Retrieve a specific recipe by its id.

* Path Parameter:

```
id: The ID of the recipe.

```

 Response

```
{
    "message": "Recipe deleted successfully"
}

```

# Error Responses
Common error responses across endpoints:

* 400 Bad Request: Invalid input.
* 401 Unauthorized: Authentication or token issues.
* 404 Not Found: Resource not found, such as missing recipe.
* 500 Internal Server Error: Server-side issues.

## Prerequisites

Before running the Recipes API, ensure you have the following installed:
* Go 1.20 or later
* Git
* (Optional) Redis for session storage
* (Optional) Docker for containerized deployment

## Installation & Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/recipes-api.git
cd recipes-api
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
Create a `.env` file (or set environment variables) with the following:
```
AUTH0_DOMAIN=your-auth0-domain.auth0.com
AUTH0_CLIENT_ID=your-client-id
JWT_SECRET=your-jwt-secret
REDIS_URL=redis://localhost:6379 (optional)
PORT=8080
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## Project Structure

The repository contains multiple implementation variants:

* **api/** - Core API implementation with basic authentication
* **auth0/** - Implementation using Auth0 for authentication
* **cookies/** - Implementation using cookie-based sessions
* **handlers/** - HTTP request handlers for recipes and authentication
* **models/** - Data models for recipes and users
* **tmp/** - Temporary build outputs

Each implementation follows the same API contract but uses different authentication mechanisms.

## Running Tests

Execute unit tests with:
```bash
go test ./...
```

For verbose output:
```bash
go test -v ./...
```

## API Documentation

The API is documented using Swagger/OpenAPI specifications. View the specification in `swagger.json`.

To view the API documentation:
1. Open the Swagger UI at `http://localhost:8080/swagger/index.html` (if configured)
2. Or use the `swagger.json` file in this repository

## Performance Testing

Use the included Apache Benchmark script for load testing:
```bash
./apache-benchmark.p
```

This will run performance tests against the API endpoints.

## Building for Production

Create a production build:
```bash
go build -o recipes-api
```

Run the compiled binary:
```bash
./recipes-api
```

## Docker Deployment

Build a Docker image:
```bash
docker build -t recipes-api:latest .
```

Run the container:
```bash
docker run -p 8080:8080 --env-file .env recipes-api:latest
```

## Troubleshooting

**JWT Token Expired**: Call the `/Refresh` endpoint with a valid token to get a new one.

**401 Unauthorized**: Ensure you've obtained a valid JWT token via `/signin` and included it in the Authorization header as `Bearer <token>`.

**Recipe Not Found**: Verify the recipe ID is correct and the recipe exists using `GET /recipes`.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add your feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. See LICENSE file for details.