# TODO List Application

## Overview

This is a TODO List application built with Go, utilizing Docker and Swagger for API documentation. The application supports user authentication, task management, and provides a RESTful API for interacting with tasks and users.

## Features

- **User Authentication:** Register, login, and manage users with JWT token-based authentication.
- **Task Management:** Create, read, update, and delete tasks.
- **Filtering and Sorting:** Filter tasks by title and status, and sort tasks by various fields.
- **API Documentation:** Swagger documentation for API endpoints.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) installed on your machine.
- [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine.

### Running the Application

1. **Clone the repository:**

   ```bash
   git clone github.com/yrss1/todo
   cd todo
   ```
2. **Build and start the application using Docker Compose:** 
   ```bash
   docker-compose up --build
    ```
   This command will build the Docker images and start the containers for the application and the PostgreSQL database.
3. **Access the API:**
   The API will be available at http://localhost:8080/api/v1.
4. **Access Swagger Documentation:**
   The Swagger UI is available at http://localhost:8080/swagger/index.html. This provides interactive documentation for the API endpoints.
   API Endpoints
   Authentication

## API Endpoints

### Authentication

- **POST /auth/login**: Login a user and receive a JWT token.
- **POST /auth/register**: Register a new user.

### Health Check

- **GET /health**: Check the health of the application.

### Tasks

- **GET /tasks**: Get all tasks with optional filtering and sorting.
- **POST /tasks**: Add a new task.
- **GET /tasks/{id}**: Get task by ID.
- **PUT /tasks/{id}**: Update task by ID.
- **DELETE /tasks/{id}**: Delete task by ID.

### Users

- **GET /users**: Get all users.
- **POST /users**: Add a new user.
- **GET /users/{id}**: Get user by ID.
- **PUT /users/{id}**: Update user by ID.
- **DELETE /users/{id}**: Delete user by ID.
- **GET /users/email**: Get user details by email.
- **GET /users/search**: Search users by name or email.

## Configuration

The application configuration is handled via environment variables. You can set the required environment variables in a `.env` file or directly in your Docker Compose configuration.

## Troubleshooting

If you encounter issues, ensure that:

- Docker and Docker Compose are properly installed and running.
- All necessary environment variables are correctly set.
- Check the application logs for any errors that might indicate the problem.

