# A Simple Template RESTful API Project by Golang Using ORM, JWT 

This is a template project for building a RESTful API using <strong>Golang</strong>, integrating <strong>GORM</strong> (ORM) for database operations, and <strong>JWT (JSON Web Tokens)</strong> for authentication. It provides a solid starting point for building scalable and secure API applications. 

## Features

<ul>
  <li>User authentication with JWT (login, register, token validation).</li>
  <li>CRUD operations with GORM.</li>
  <li>Middleware for securing endpoints.</li>
  <li>Configuration via environment variables.</li>
  <li>Modular and clean project structure.</li>
  <li>Support for PostgreSQL, MySQL, and SQLite (via GORM).</li>
</ul>

## Prerequisites

<ul>
  <li>Go 1.19 or newer</li>
  <li>A database (PostgreSQL, MySQL, or SQLite)</li>
  <li><a href="https://www.postman.com/" target="_blank">Postman</a> (optional, for API testing)</li>
</ul>

## Installation

* Clone the repository:
   ``` bash
    git clone https://github.com/tnqbao/goang_template.git
    cd golang_template
   ```
* Setup your module: 
  ``` bash
   go mod edit -module=your-link-github-repo 
  ```
* Install dependencies:
  ``` bash
    go mod tidy 
  ``` 
  
* Set up environment variables:
    * Create a `.env` file in the project root and configure the following variables:
  ```dotenv
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_HOST=localhost or your db host
    DB_PORT= your db port
    ## your other env
    ```
* Run database migrations:
    ``` bash
     go run main.go migrate
   ```

* Start the server:
    ``` bash 
    go run main.go
    ```
 
  <li>Access the API at: <a href="http://localhost:8080" target="_blank">http://localhost:8080</a></li>

## Project Structure
   ``` 
   ├── config/         # Configuration files
   ├── api/            # Handlers for API routes
   ├── middlewares/    # Middleware (e.g., JWT authentication)
   ├── models/         # GORM models
   ├── routes/         # API route definitions
   ├── utils/          # Utility functions
   ├── main.go         # Entry point of the application
   └── go.mod          # Go module file
   ```



<h2>Future Improvements</h2>
<ul>
  <li>Add unit and integration tests.</li>
  <li>Implement role-based access control (RBAC).</li>
  <li>Add API versioning.</li>
  <li>Improve error handling and logging.</li>
</ul>

<h2>License</h2>
<p>This project is licensed under the MIT License. See the <a href="LICENSE" target="_blank">LICENSE</a> file for details.</p>
