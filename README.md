# gau_phim_backend

![gau_phim_backend](https://img.shields.io/badge/Service-gau_phim_backend-blue)

`gau_phim_backend` is a backend service built using **Golang**. This service provides APIs for managing movie-related data, user authentication, and movie history tracking. It integrates **Redis** for caching, **PostgreSQL** for data storage, **MeiliSearch** for advanced movie search functionalities, and features an **AI-powered chatbot** for movie search assistance.

## Features

- 🔒 **JWT Authentication**: Secure endpoints with JWT-based authentication.
- 🎬 **Movie Management**: APIs to manage movie data such as creating, updating, and retrieving movie details.
- 📜 **User History**: Track users' movie history with the ability to add, update, or delete history records.
- 🔍 **MeiliSearch Integration**: Fast, full-text search for movies.
- 🧠 **AI Chatbot**: An intelligent chatbot to help users find movies based on natural language queries.
- 🧑‍💻 **Redis Cache**: Caching for frequent queries to improve performance.
- 🗄️ **PostgreSQL Database**: Data persistence for movie records and user history.

## Tech Stack

- **Golang**: The service is built using Golang for high performance.
- **JWT**: JSON Web Token (JWT) for secure user authentication.
- **GORM**: ORM (Object-Relational Mapping) for interacting with PostgreSQL.
- **Redis**: In-memory data store used for caching and session management.
- **PostgreSQL**: Relational database for storing movie data and user history.
- **MeiliSearch**: A powerful search engine used for indexing and searching movies.
- **AI Chatbot**: A conversational AI to help users find movies based on queries.

## Installation

### Prerequisites

Before running the service, make sure you have the following installed:

- [Golang](https://golang.org/dl/) (v1.18+)
- [Docker](https://www.docker.com/get-started)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download)
- [MeiliSearch](https://www.meilisearch.com/docs/)

### Step-by-Step Installation

1. **Clone the repository**:

    ```bash
    git clone https://github.com/tnqbao/gau_phim_backend.git
    cd gau_phim_backend
    ```

2. **Set up environment variables**:

   Create a `.env` file in the root directory and populate it with the following values:

    ```bash
    JWT_SECRET_KEY=your_jwt_secret_key
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_DB=gau_phim
    POSTGRES_USER=your_postgres_user
    POSTGRES_PASSWORD=your_postgres_password
    REDIS_HOST=localhost
    REDIS_PORT=6379
    MEILISEARCH_HOST=http://localhost:7700
    ```

3. **Install dependencies**:

   Install Go dependencies using the Go module system:

    ```bash
    go mod tidy
    ```

4. **Run the Docker containers (optional)**:

   If you're using Docker, you can set up a PostgreSQL and Redis container by running:

    ```bash
    docker-compose up
    ```

5. **Run the service**:

   Start the backend server with:

    ```bash
    go run main.go
    ```

## API Test Documentation

For API tests, please refer to our [API Test Documentation](https://nzs2rc26yh.apidog.io).

## Database Structure

- **Movie**: Stores movie details.
- **History**: Stores the user's watch history.
- **Search Index**: Managed by MeiliSearch for fast searching.

## Caching

The backend utilizes **Redis** to cache frequently accessed movie data and search results. This significantly improves response times for high-demand endpoints.

## Running Tests

To run the tests, make sure the environment is set up correctly and run:

```bash
go test ./...
