package ginrequests

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNormalizeRequests_SingleSlice(t *testing.T) {
	handler := func(c *gin.Context) {}
	requests := BuildRequests("GET", "/path", handler)

	result := NormalizeRequests(requests)

	if len(result) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(result))
	}
}

func TestNormalizeRequests_MultipleSlices(t *testing.T) {
	handler := func(c *gin.Context) {}

	req1 := BuildRequests("GET", "/path1", handler)
	req2 := BuildRequests("POST", "/path2", handler)
	req3 := BuildRequests("DELETE", "/path3", handler)

	result := NormalizeRequests(req1, req2, req3)

	if len(result) != 3 {
		t.Fatalf("Expected 3 requests, got %d", len(result))
	}
}

func TestNormalizeRequests_EmptySlice(t *testing.T) {
	result := NormalizeRequests()

	if len(result) != 0 {
		t.Fatalf("Expected 0 requests for empty input, got %d", len(result))
	}
}

func TestNormalizeRequests_EmptyAndNonEmptySlices(t *testing.T) {
	handler := func(c *gin.Context) {}

	req1 := BuildRequests("GET", "/path1", handler)
	emptyReq := []Request{}
	req2 := BuildRequests("POST", "/path2", handler)

	result := NormalizeRequests(req1, emptyReq, req2)

	if len(result) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(result))
	}
}

func TestNormalizeRequests_PreservesOrder(t *testing.T) {
	handler1 := func(c *gin.Context) { c.String(200, "1") }
	handler2 := func(c *gin.Context) { c.String(200, "2") }
	handler3 := func(c *gin.Context) { c.String(200, "3") }

	req1 := BuildRequests("GET", "/first", handler1)
	req2 := BuildRequests("GET", "/second", handler2)
	req3 := BuildRequests("GET", "/third", handler3)

	result := NormalizeRequests(req1, req2, req3)

	if result[0].Path != "/first" {
		t.Errorf("Expected first path '/first', got '%s'", result[0].Path)
	}
	if result[1].Path != "/second" {
		t.Errorf("Expected second path '/second', got '%s'", result[1].Path)
	}
	if result[2].Path != "/third" {
		t.Errorf("Expected third path '/third', got '%s'", result[2].Path)
	}
}

func TestNormalizeRequests_MultipleHandlersSamePath(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}

	req1 := BuildRequests("POST", "/api", handler1, handler2)
	req2 := BuildRequests("PUT", "/api/id", handler1)

	result := NormalizeRequests(req1, req2)

	if len(result) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(result))
	}
	if len(result[0].Handler) != 2 {
		t.Fatalf("First request should have 2 handlers, got %d", len(result[0].Handler))
	}
}

func TestNormalizeRequests_PreservesRequestProperties(t *testing.T) {
	handler := func(c *gin.Context) {}

	req := BuildRequests("DELETE", "/resource/:id", handler)
	result := NormalizeRequests(req)

	if result[0].Method != "DELETE" {
		t.Errorf("Expected method 'DELETE', got '%s'", result[0].Method)
	}
	if result[0].Path != "/resource/:id" {
		t.Errorf("Expected path '/resource/:id', got '%s'", result[0].Path)
	}
	if len(result[0].Handler) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(result[0].Handler))
	}
}

func TestNormalizeRequests_ManySlices(t *testing.T) {
	handler := func(c *gin.Context) {}

	slices := [][]Request{}
	expectedTotal := 0

	for i := 0; i < 10; i++ {
		req := BuildRequests("GET", "/path"+string(rune(i)), handler)
		slices = append(slices, req)
		expectedTotal += len(req)
	}

	result := NormalizeRequests(slices...)

	if len(result) != expectedTotal {
		t.Fatalf("Expected %d requests, got %d", expectedTotal, len(result))
	}
}

func TestNormalizeRequests_ComplexMixture(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}
	handler3 := func(c *gin.Context) {}

	req1 := BuildRequests("GET",
		"/users", handler1,
		"/users/:id", handler2,
	)

	req2 := BuildRequests("POST",
		"/users", handler1, handler3,
	)

	req3 := []Request{} // Empty slice

	req4 := BuildRequests("DELETE", "/users/:id", handler1)

	result := NormalizeRequests(req1, req2, req3, req4)

	if len(result) != 4 {
		t.Fatalf("Expected 4 requests total, got %d", len(result))
	}
}

func TestGetTotalRequestsLength_EmptySlices(t *testing.T) {
	slices := [][]Request{}
	total := getTotalRequestsLength(&slices)

	if total != 0 {
		t.Fatalf("Expected 0 total requests, got %d", total)
	}
}

func TestGetTotalRequestsLength_MultipleSlices(t *testing.T) {
	handler := func(c *gin.Context) {}

	req1 := BuildRequests("GET", "/path1", handler)
	req2 := BuildRequests("GET", "/path2", handler)
	req3 := BuildRequests("GET", "/path3", handler)

	slices := [][]Request{req1, req2, req3}
	total := getTotalRequestsLength(&slices)

	if total != 3 {
		t.Fatalf("Expected 3 total requests, got %d", total)
	}
}

func TestGetTotalRequestsLength_WithEmptySlices(t *testing.T) {
	handler := func(c *gin.Context) {}

	req1 := BuildRequests("GET", "/path1", handler)
	emptyReq := []Request{}
	req2 := BuildRequests("GET", "/path2", handler)

	slices := [][]Request{req1, emptyReq, req2}
	total := getTotalRequestsLength(&slices)

	if total != 2 {
		t.Fatalf("Expected 2 total requests, got %d", total)
	}
}
