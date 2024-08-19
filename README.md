# go-quick-url

## Overview
**go-quick-url** is a URL shortening service built with Go, ObjectBox, and Fiber. This project aims to provide a simple and efficient way to shorten URLs and manage them.

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
    git clone https://github.com/akhiltn/go-quick-url.git
    cd go-quick-url
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

- **POST /shorten**
    - Request: `{ "url": "https://example.com" }`
    - Response: `{ "shortened_url": "http://localhost:3000/abc123" }`

- **GET /:shortened_url**
    - Redirects to the original URL.

### Example

To shorten a URL:
```sh
curl -X POST http://localhost:3000/shorten -d '{"url":"https://example.com"}' -H "Content-Type: application/json"
