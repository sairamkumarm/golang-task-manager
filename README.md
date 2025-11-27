# Task Manager API

A lightweight task management backend built with Go, Gin, SQLx, PostgreSQL, Redis and JWT authentication.
Includes user registration and login, task CRUD operations, optional public shareable slugs, search, pagination, due date filtering, rate limiting and a full OpenAPI 3.1 specification.

## Features

* User authentication using JWT
* Create, update, delete and fetch tasks
* Public task sharing through slugs
* Search by keyword and status
* Filter by exact due date
* Pagination included by default
* Redis based rate limiting
* SQLx repository layer
* Fully containerized with Docker Compose
* OpenAPI 3.1 compliant API spec

## Tech Stack

* Go 1.21
* Gin
* SQLx
* PostgreSQL 15
* Redis 7
* JWT
* Docker and Docker Compose

---

## Getting Started

### Prerequisites

* Docker
* Docker Compose

### Environment Variables

`.env` (example has same credentials as docker images, feel free to use the same):

```
# Server
PORT=8080
JWT_SECRET=7627cf4ae27fa73abff9613a74f701d6a3efd41762f1649a7acf3cecaee6c962
JWT_EXPIRATION_HOURS=72

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=resu
DB_PASSWORD=drowssap
DB_NAME=task-manager-db

# Redis (rate limiting)
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=drowssap
RATE_LIMIT=10        # requests per minute per IP

# Optional
LOG_LEVEL=info
```

### Running the application

```
docker compose up --build
```

App runs on:

```
http://localhost:8080
```

---

## Project Structure

```
/cmd
    main.go
/internal
    /handler
    /service
    /repository
    /middleware
    /models
    /dtos
    /utils
    /logger
    /migrations
```
---

## API Documentation

The API follows OpenAPI 3.1.0.  
This is a compact, readable summary of the routes.  

### Authentication

#### POST /api/register

Registers a new user.  
Request body uses `RegisterRequest`.  
Returns `201 Created`.  

#### POST /api/login

Authenticates a user and returns a JWT.  
Request body uses `LoginRequest`.  

Response:

```json
{
  "token": "jwt-token"
}
```

Include it in header for protected routes:

```
Authorization: Bearer <token>
```

---

## Tasks

### POST /api/tasks

Create a task.  
Requires authentication.  

Body: `CreateTaskRequest`

Returns `201 Created` and a `TaskResponse`.

---

### GET /api/tasks

List tasks for the logged in user.  
Supports filters:  

| Query    | Type      | Description                     |
| -------- | --------- | ------------------------------- |
| status   | string    | pending, in_progress, completed |
| keyword  | string    | search in title and description |
| due_date | date-time | exact date match, ignores time  |
| page     | int       | default 1                       |
| limit    | int       | default 10                      |

Response body: `ListTasksResponse`

---

### GET /api/tasks/{id}

Fetch a single task by UUID.  
Requires authentication.  

---

### PUT /api/tasks/{id}

Update a task.  
Requires authentication.   
Body: `UpdateTaskRequest`

---

### DELETE /api/tasks/{id}

Delete a task.  
Requires authentication.  
Returns `204 No Content`.  

---

## Public Tasks

### GET /api/public/{slug}

Fetch a publicly shared task.
No authentication.

Response: `TaskResponse`

---

## Data Schemas

### RegisterRequest

```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

### LoginRequest

```json
{
  "email": "string",
  "password": "string"
}
```

### CreateTaskRequest

```json
{
  "title": "string",
  "description": "string",
  "status": "pending | in_progress | completed",
  "due_date": "date-time or null",
  "is_public": true | false
}
```

### UpdateTaskRequest

All fields optional.

```json
{
  "title": "string",
  "description": "string",
  "status": "pending | in_progress | completed",
  "due_date": "date-time or null",
  "is_public": true | false
}
```

### TaskResponse

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "string",
  "description": "string",
  "status": "pending | in_progress | completed",
  "due_date": "date-time or null",
  "is_public": boolean,
  "public_slug": "string or null",
  "created_at": "date-time",
  "updated_at": "date-time"
}
```

### ListTasksResponse

```json
{
  "tasks": [ TaskResponse ],
  "page": number,
  "limit": number,
  "total": number
}
```

---

### Error Responses

All errors follow the format:

```json
{
  "error": "message"
}
```

---

If you want this in a canvas file or want it split into README.md plus a separate docs file, say the word.
