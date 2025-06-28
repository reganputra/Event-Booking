# Event Booking REST API

A RESTful API built with Go and Gin framework for managing events and user registrations. This API allows users to create, view, update, and delete events, as well as register for events. Now supports user roles (admin and user) for enhanced access control.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [User Management](#user-management)
  - [Event Management](#event-management)
  - [Event Registration](#event-registration)
  - [Admin Endpoints](#admin-endpoints)
- [Authentication & Authorization](#authentication--authorization)

## Features

- User registration and authentication with JWT
- User roles: `user` and `admin`
- Admin-only endpoints for user management
- CRUD operations for events
- Event registration functionality
- Protected routes with middleware authentication and role-based authorization
- SQLite database for data storage

## Technologies

- Go (Golang)
- Gin Web Framework
- SQLite
- JWT for authentication

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-rest-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. The server will start on port 3000 by default.

## API Endpoints

### Health Check

- **GET /healthcheck** - Check if the server is running

### User Management

- **POST /users/register** - Register a new user
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
  - Response:
    ```json
    {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "role": "user"
      }
    }
    ```

- **POST /users/login** - Login a user
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
  - Response:
    ```json
    {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "role": "user"
      },
      "token": "jwt-token-here"
    }
    ```

### Event Management

- **GET /events** - Get all events (public)
  - Response: Array of event objects

- **GET /events/:id** - Get a specific event by ID (public)
  - Response: Event object

- **POST /events** - Create a new event (protected, any authenticated user)
  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "name": "New Event",
      "description": "This is a new event description",
      "location": "123 Event St, Event City, EC 12345",
      "date": "2023-12-01T15:00:00Z"
    }
    ```
  - Response:
    ```json
    {
      "message": "Event created successfully!",
      "event": {
        "id": 1,
        "name": "New Event",
        "description": "This is a new event description",
        "location": "123 Event St, Event City, EC 12345",
        "date": "2023-12-01T15:00:00Z",
        "userIds": 1
      }
    }
    ```

- **PUT /events/:id** - Update an event (protected, owner or admin)
  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "name": "Updated Event",
      "description": "This is an updated event description",
      "location": "456 Updated St, Updated City, UC 67890",
      "date": "2023-12-15T16:00:00Z"
    }
    ```
  - Response:
    ```json
    {
      "message": "Event updated successfully!",
      "event": {
        "id": 1,
        "name": "Updated Event",
        "description": "This is an updated event description",
        "location": "456 Updated St, Updated City, UC 67890",
        "date": "2023-12-15T16:00:00Z",
        "userIds": 1
      }
    }
    ```

- **DELETE /events/:id** - Delete an event (protected, owner or admin)
  - Headers: `Authorization: Bearer <token>`
  - Response:
    ```json
    {
      "message": "Event deleted successfully!"
    }
    ```

### Event Registration

- **POST /events/:id/register** - Register for an event (protected)
  - Headers: `Authorization: Bearer <token>`
  - Response:
    ```json
    {
      "message": "Successfully registered for the event"
    }
    ```

- **DELETE /events/:id/register** - Cancel registration for an event (protected)
  - Headers: `Authorization: Bearer <token>`
  - Response:
    ```json
    {
      "message": "Successfully cancelled event registration"
    }
    ```

### Admin Endpoints

Admin endpoints require the user to have the `admin` role. Use the JWT token of an admin user in the `Authorization` header.

- **GET /admin/users** - Get all users
- **GET /admin/users/:id** - Get user by ID
- **PUT /admin/users/:id** - Update a user (role and email)
- **DELETE /admin/users/:id** - Delete a user

#### Example Admin Request

```http
GET /admin/users
Authorization: Bearer <admin-jwt-token>
```

## Authentication & Authorization

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. Register a user or login to get a JWT token
2. Include the token in the Authorization header for protected requests:
   - `Authorization: Bearer <token>` or
   - `Authorization: <token>`

### User Roles
- **user**: Can register/login, view and manage their own events, register for events.
- **admin**: Has all user permissions plus access to admin endpoints for managing users.

Role-based access is enforced using middleware. Admin endpoints are only accessible to users with the `admin` role.
