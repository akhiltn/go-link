# go-link

## Overview
**go-link** is a URL shortening service built with Go, BoltDB, and Fiber. This project aims to provide a simple and efficient way to shorten URLs and manage them.

## Features
- URL shortening
- URL redirection
- URL management
- RESTful API
- Swagger API documentation

## Technologies Used
- **Go**: The main programming language used for the project.
- **BoltDB**: A high-performance NoSQL database for Go.
- **Fiber**: An Express-inspired web framework for Go.
- **Swagger**: For API documentation.

## Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/akhiltn/go-link.git
    cd go-link
    ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

3. **Run the application:**
    ```sh
    go run main.go
    ```

## Usage

### API Endpoints

- **POST /
    - Request: `{ "url": "https://example.com", "short": "example" }`

- **GET /:shortened_url**
    - Redirects to the original URL.

### Example

To shorten a URL:
```sh
curl -X POST http://localhost:3000/ -d '{"url":"https://example.com", "short": "example" }' -H "Content-Type: application/json"
