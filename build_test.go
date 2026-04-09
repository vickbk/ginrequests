package ginrequests

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBuildRequests_SinglePath_SingleHandler(t *testing.T) {
	handler := func(c *gin.Context) {}

	requests := BuildRequests("GET", "/path", handler)

	if len(requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requests))
	}
	if requests[0].Path != "/path" {
		t.Fatalf("Expected path '/path', got '%s'", requests[0].Path)
	}
	if requests[0].Method != "GET" {
		t.Fatalf("Expected method 'GET', got '%s'", requests[0].Method)
	}
	if len(requests[0].Handler) != 1 {
		t.Fatalf("Expected 1 handler, got %d", len(requests[0].Handler))
	}
}

func TestBuildRequests_MultiplePaths(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}

	requests := BuildRequests("GET",
		"/path1", handler1,
		"/path2", handler2,
	)

	if len(requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(requests))
	}
	if requests[0].Path != "/path1" {
		t.Fatal("First request path mismatch")
	}
	if requests[1].Path != "/path2" {
		t.Fatal("Second request path mismatch")
	}
}

func TestBuildRequests_MultipleHandlersSamePath(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}

	requests := BuildRequests("POST",
		"/path", handler1, handler2,
	)

	if len(requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requests))
	}
	if len(requests[0].Handler) != 2 {
		t.Fatalf("Expected 2 handlers, got %d", len(requests[0].Handler))
	}
}

func TestBuildRequests_MixedPathsAndHandlers(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}
	handler3 := func(c *gin.Context) {}

	requests := BuildRequests("DELETE",
		"/path1", handler1,
		"/path2", handler2, handler3,
	)

	if len(requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(requests))
	}
	if len(requests[0].Handler) != 1 {
		t.Fatalf("First request should have 1 handler, got %d", len(requests[0].Handler))
	}
	if len(requests[1].Handler) != 2 {
		t.Fatalf("Second request should have 2 handlers, got %d", len(requests[1].Handler))
	}
}

func TestBuildRequests_ValidatesRequests(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for invalid request configuration")
		}
	}()

	// Create request with empty path
	requests := BuildRequests("GET", "")
	_ = requests
}

func TestBuildRequests_MethodPreserved(t *testing.T) {
	handler := func(c *gin.Context) {}

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}

	for _, method := range methods {
		requests := BuildRequests(method, "/path", handler)
		if requests[0].Method != method {
			t.Errorf("Expected method '%s', got '%s'", method, requests[0].Method)
		}
	}
}

func TestBuildRequests_ComplexScenario(t *testing.T) {
	h1 := func(c *gin.Context) {}
	h2 := func(c *gin.Context) {}
	h3 := func(c *gin.Context) {}
	h4 := func(c *gin.Context) {}

	requests := BuildRequests("GET",
		"/users", h1, h2,
		"/posts", h3,
		"/comments", h4,
	)

	if len(requests) != 3 {
		t.Fatalf("Expected 3 requests, got %d", len(requests))
	}
	if len(requests[0].Handler) != 2 {
		t.Fatal("First request should have 2 handlers")
	}
	if len(requests[1].Handler) != 1 {
		t.Fatal("Second request should have 1 handler")
	}
	if len(requests[2].Handler) != 1 {
		t.Fatal("Third request should have 1 handler")
	}
}

func TestBuildRequests_EmptyPaths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for empty path")
		}
	}()

	handler := func(c *gin.Context) {}
	requests := BuildRequests("GET", "", handler)
	_ = requests
}

func TestBuildRequests_NoHandlers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for path with no handlers")
		}
	}()

	requests := BuildRequests("GET", "/path")
	_ = requests
}

func TestBuildRequests_ParameterizedPaths(t *testing.T) {
	handler := func(c *gin.Context) {}

	requests := BuildRequests("GET",
		"/users/:id", handler,
		"/posts/:id/comments/:commentId", handler,
	)

	if len(requests) != 2 {
		t.Fatal("Expected 2 requests")
	}
	if requests[0].Path != "/users/:id" {
		t.Errorf("Expected path '/users/:id', got '%s'", requests[0].Path)
	}
	if requests[1].Path != "/posts/:id/comments/:commentId" {
		t.Errorf("Expected path '/posts/:id/comments/:commentId', got '%s'", requests[1].Path)
	}
}

func TestIsHandlerFunc_ValidHandler(t *testing.T) {
	handler := func(c *gin.Context) {}
	if !isHandlerFunc(handler) {
		t.Fatal("Should recognize valid handler func")
	}
}

func TestIsHandlerFunc_InvalidInputs(t *testing.T) {
	tests := []struct {
		input interface{}
		name  string
	}{
		{"string", "string"},
		{123, "integer"},
		{12.34, "float"},
		{[]string{}, "slice"},
		{map[string]int{}, "map"},
		{func() {}, "no-param func"},
	}

	for _, test := range tests {
		if isHandlerFunc(test.input) {
			t.Errorf("Should not recognize %s as handler func", test.name)
		}
	}
}
