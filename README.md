# http-go
This is a simple HTTP server written in Go, created as a learning project. The project implements the basic principles of HTTP request handling, routing, and response generation, with the goal of understanding how HTTP servers function at a low level. This is not a production-ready HTTP server and is intended to provide insight into how HTTP works in Go.

## Key Features

- **Custom HTTP Server**: A basic HTTP server that can handle GET requests and return responses.
- **Request Handling**: Supports custom handling of HTTP requests with routing to different pre-defined endpoints.
- **Minimal Dependencies**: The server is built using Go's standard library with no external dependencies

## Installation

### Running the Project

Once you have cloned the project, you can run the server with the following command:
```bash
go run .
```

This will start the HTTP server, which will listen for incoming HTTP requests. By default, it listens on localhost:8080. You can access the server via a browser or use a tool like curl to send requests.

For example, you can test the server by visiting http://localhost:8080/ in your web browser, or by running a curl request.
