# gin-requests

A Go package that simplifies building and managing HTTP request handlers for the Gin web framework. It provides a structured way to define routes with their handlers, making it easier to organize and maintain your Gin applications.

## Features

- **Structured Request Building**: Define HTTP requests with paths and handlers in a clean, readable format
- **Multiple Handlers Support**: Attach multiple middleware or handlers to the same route
- **Request Normalization**: Combine multiple request slices into a single collection
- **Type Safety**: Strongly typed request structures with validation
- **Gin Integration**: Seamlessly works with the Gin web framework

## Installation

```bash
go get github.com/vickbk/ginrequests
```

## Usage

### Basic Example

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/vickbk/ginrequests"
)

func main() {
    r := gin.Default()

    // Define handlers
    getAlbumsHandler := func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"albums": []string{"album1", "album2"}})
    }

    createAlbumHandler := func(c *gin.Context) {
        c.JSON(http.StatusCreated, gin.H{"message": "Album created"})
    }

    // Build requests using the package
    requests := ginrequests.BuildRequests("GET",
        "/albums", getAlbumsHandler,
        "/health", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        },
    )

    postRequests := ginrequests.BuildRequests("POST",
        "/albums", createAlbumHandler,
    )

    // Normalize all requests into a single slice
    allRequests := ginrequests.NormalizeRequests(requests, postRequests)

    // Register routes with Gin
    allRequests.AddRoutes(r)

    r.Run(":8080")
}
```

### Advanced Example with Multiple Handlers

```go
// Define middleware and handlers
authMiddleware := func(c *gin.Context) {
    // Authentication logic
    c.Next()
}

validateAlbumHandler := func(c *gin.Context) {
    // Validation logic
    c.Next()
}

createAlbumHandler := func(c *gin.Context) {
    c.JSON(http.StatusCreated, gin.H{"message": "Album created"})
}

// Build requests with multiple handlers per route
albumRequests := ginrequests.BuildRequests("POST",
    "/albums", authMiddleware, validateAlbumHandler, createAlbumHandler,
)
```

## API Reference

### Types

#### `Request`

Represents an HTTP request configuration.

```go
type Request struct {
    Handler      []gin.HandlerFunc
    Method, Path string
}
```

#### `RequestList`

A slice of `Request` objects.

```go
type RequestList []Request
```

### Functions

#### `BuildRequests(method string, rqs ...any) RequestList`

Constructs a slice of `Request` objects from the provided HTTP method and variadic arguments.

**Parameters:**

- `method`: HTTP method (e.g., "GET", "POST", "PUT", "DELETE")
- `rqs`: Alternating sequence of path strings and handler functions

**Returns:** A `RequestList` containing the built requests

**Example:**

```go
requests := BuildRequests("GET",
    "/users", getUsersHandler,
    "/users/:id", getUserByIdHandler,
)
```

#### `NormalizeRequests(rs ...[]Request) []Request`

Flattens multiple slices of `Request` objects into a single combined slice.

**Parameters:**

- `rs`: Variable number of `[]Request` slices

**Returns:** A single slice containing all requests

**Example:**

```go
allRequests := NormalizeRequests(getRequests, postRequests, deleteRequests)
```

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source. Please check the license terms if any are specified.
