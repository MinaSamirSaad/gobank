# Gobank - A Simple Bank API

## Overview

**Gobank** is a simple RESTful API built in **Go** using the **Mux** router. It simulates basic banking operations such as account creation, retrieval by account ID, account deletion, and user authentication with JWT tokens.

### Key Features:
- **Create an account**: Register a new account with a specified balance.
- **Get account details by ID**: Retrieve account information by account ID.
- **Delete an account**: Remove an account from the system.
- **Sign in & Sign out**: Authenticate users with JWT and allow sign-out functionality.

The API uses **PostgreSQL** for persistent storage and **Docker** for containerization of both the application and the database.

## Technologies Used
- **Go (Golang)** for building the backend REST API.
- **Mux** as the HTTP router for routing requests.
- **JWT (JSON Web Tokens)** for secure authentication.
- **PostgreSQL** as the database to store user data and account information.
- **Docker** for containerizing both the application and the PostgreSQL database.

## Prerequisites

Before running this project locally, ensure that the following are installed:

- **Go (Golang)**: Follow the [Go installation guide](https://golang.org/doc/install) to install the Go programming language.
- **Docker**: For running PostgreSQL in a container, follow the [Docker installation guide](https://www.docker.com/get-started).
- **Make** (Optional but recommended): This tool simplifies the process of building and running the application. Install Make from [here](https://www.gnu.org/software/make/).

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/MinaSamirSaad/gobank.git
cd gobank
```

### 3. Run the Application

Use the following command to build and run the application:

```bash
make run
```

This command will start the application and the PostgreSQL database in Docker containers.

### 4. Access the API

Once the application is running, you can access the API at `http://localhost:3000`. Use tools like `curl` or `Postman` to interact with the endpoints.

## API Endpoints

Here are some of the key endpoints available in the Gobank API:

- **POST /account**: Create a new account.
- **GET /account/{id}**: Retrieve account details by account ID.
- **DELETE /account/{id}**: Delete an account by account ID.
- **POST /login**: Sign in a user and receive a JWT token.
- **POST /logout**: Sign out a user.

## Conclusion

Gobank is a simple yet powerful API for managing basic banking operations. With its use of Go, PostgreSQL, and Docker.

Happy coding!
