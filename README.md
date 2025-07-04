# Event Booking REST API

A RESTful API built with Go and Gin framework for managing events and user registrations. This API allows users to create, view, update, and delete events, as well as register for events. Now supports user roles (admin and user) for enhanced access control.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Installation](#installation)
  - [Using Docker (Recommended)](#using-docker-recommended)
  - [Local Installation](#local-installation)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [User Management](#user-management)
  - [Event Management](#event-management)
  - [Event Registration](#event-registration)
  - [Event Reviews](#event-reviews)
  - [Event Waitlist](#event-waitlist)
  - [Admin Endpoints](#admin-endpoints)
- [Authentication & Authorization](#authentication--authorization)

## Features

- User registration and authentication with JWT
- User roles: `user` and `admin`
- Admin-only endpoints for user management
- CRUD operations for events
- Event registration functionality
- Event search and filtering (by keyword, date range)
- Event reviews and ratings
- Waitlist system for full events
- Protected routes with middleware authentication and role-based authorization
- PostgreSQL database for data storage
- Docker support for easy setup and deployment

## Technologies

- Go (Golang)
- Gin Web Framework
- PostgreSQL
- Docker
- JWT for authentication

## Installation

### Using Docker (Recommended)

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd go-rest-api
    ```

2.  **Run the application with Docker Compose:**
    This command will build the Go application and start the PostgreSQL database container.

    ```bash
    docker-compose up --build
    ```

3.  The server will start on port `3000`, and the PostgreSQL database will be available on port `5432`.

### Local Installation

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd go-rest-api
    ```

2.  **Install dependencies:**

    ```bash
    go mod download
    ```

3.  **Set up Environment Variables:**
    Create a `.env` file in the root directory and add the following environment variable for the database connection.

    ```
    DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
    ```

    If you are not using the default credentials, update the connection string accordingly.

4.  **Run the application:**

    ```bash
    go run main.go
    ```

5.  The server will start on port `3000` by default.

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

- **GET /events/category/:category** - Get events by category (public)

  - Response: Array of event objects

- **GET /events/search** - Search events by keyword, start date, or end date (public)
  - Query Parameters:
    - `keyword` (string, optional): Search term for event name or description.
    - `startDate` (string, optional, format: `YYYY-MM-DD`): Filter events starting on or after this date.
    - `endDate` (string, optional, format: `YYYY-MM-DD`): Filter events ending on or before this date.
  - Example: `/events/search?keyword=Workshop&startDate=2024-03-01`
  - Response: Array of event objects

- **POST /events** - Create a new event (protected, any authenticated user)

  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "name": "New Event",
      "description": "This is a new event description",
      "location": "123 Event St, Event City, EC 12345",
      "date": "2023-12-01T15:00:00Z",
      "category": "Tech",
      "capacity": 50  // Optional: Maximum number of attendees. 0 or omitted for unlimited.
    }
    ```
  - Response:
    ```json
    {
      "message": "Event created successfully!",
      "event": { ... }
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
      "date": "2023-12-15T16:00:00Z",
      "category": "Health"
    }
    ```
  - Response:
    ```json
    {
      "message": "Event updated successfully!",
      "event": { ... }
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
  - Note: If the event is full and has a capacity set, this might return a 202 Accepted with a message like:
    ```json
    {
      "message": "event is full, user added to waitlist"
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

- **GET /events/registered** - Get all events a user is registered for (protected)
  - Headers: `Authorization: Bearer <token>`
  - Response: Array of event objects

### Event Reviews

- **POST /events/:id/reviews** - Create a review for an event (protected)
  - User must be authenticated. (Future enhancement: ensure user was registered for the event).
  - Headers: `Authorization: Bearer <token>`
  - Request body:
    ```json
    {
      "rating": 5, // Integer between 1 and 5
      "comment": "This was an amazing event!"
    }
    ```
  - Response (201 Created):
    ```json
    {
      "message": "Review created successfully",
      "review": {
        "id": 1,
        "event_id": 123,
        "user_id": 1,
        "rating": 5,
        "comment": "This was an amazing event!",
        "created_at": "2024-03-15T10:00:00Z"
      }
    }
    ```
  - Response (409 Conflict if already reviewed):
    ```json
    {
      "error": "you have already reviewed this event"
    }
    ```

- **GET /events/:id/reviews** - Get all reviews for a specific event (public)
  - Response: Array of review objects. Each event object returned from `/events` or `/events/:id` will also now include an `average_rating` field.

### Event Waitlist

- **POST /events/:id/waitlist** - Join the waitlist for a full event (protected)
  - User must be authenticated.
  - This endpoint should typically be called if `POST /events/:id/register` indicates the event is full and the user was added to the waitlist, or if a user explicitly wants to join a known full event's waitlist.
  - Headers: `Authorization: Bearer <token>`
  - Response (201 Created):
    ```json
    {
      "message": "Successfully joined the waitlist",
      "waitlist_entry": {
        "id": 1,
        "event_id": 123,
        "user_id": 1,
        "created_at": "2024-03-15T11:00:00Z"
      }
    }
    ```
  - Response (409 Conflict if various conditions not met, e.g., event not full, already registered, already on waitlist):
    ```json
    {
      "error": "Specific error message like 'event is not full, cannot join waitlist'"
    }
    ```

- **DELETE /events/:id/waitlist** - Leave the waitlist for an event (protected)
  - User must be authenticated.
  - Headers: `Authorization: Bearer <token>`
  - Response (200 OK):
    ```json
    {
      "message": "Successfully left the waitlist"
    }
    ```
  - Response (404 Not Found if user not on waitlist):
    ```json
    {
      "error": "user is not on the waitlist for this event"
    }
    ```

### Admin Endpoints

Admin endpoints require the user to have the `admin` role. Use the JWT token of an admin user in the `Authorization` header.

- **GET /admin/users** - Get all users
- **GET /admin/users/:id** - Get user by ID
- **PUT /admin/users/:id** - Update a user (role and email)
- **DELETE /admin/users/:id** - Delete a user
- **GET /admin/events/:id/waitlist** - Get the waitlist for a specific event (admin)
  - Headers: `Authorization: Bearer <admin-jwt-token>`
  - Response: Array of waitlist entry objects.

#### Example Admin Request

```http
GET /admin/users
Authorization: Bearer <admin-jwt-token>
```

## Authentication & Authorization

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1.  Register a user or login to get a JWT token
2.  Include the token in the Authorization header for protected requests:
    - `Authorization: Bearer <token>`

### User Roles

- **user**: Can register/login, view and manage their own events, register for events.
- **admin**: Has all user permissions plus access to admin endpoints for managing users.

Role-based access is enforced using middleware. Admin endpoints are only accessible to users with the `admin` role.
