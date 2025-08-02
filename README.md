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