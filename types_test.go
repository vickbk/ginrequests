package ginrequests

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequest_Validate_ValidRequest(t *testing.T) {
	handler := func(c *gin.Context) {}
	req := Request{
		Path:    "/valid",
		Method:  "GET",
		Handler: []gin.HandlerFunc{handler},
	}

	err := req.validate()
	if err != nil {
		t.Fatalf("Valid request should not produce error: %v", err)
	}
}

func TestRequest_Validate_EmptyPath(t *testing.T) {
	handler := func(c *gin.Context) {}
	req := Request{
		Path:    "",
		Method:  "GET",
		Handler: []gin.HandlerFunc{handler},
	}

	err := req.validate()
	if err == nil {
		t.Fatal("Empty path should produce error")
	}
	if err.Error() != "request path cannot be empty" {
		t.Fatalf("Expected 'request path cannot be empty' error, got: %v", err)
	}
}

func TestRequest_Validate_NoHandlers(t *testing.T) {
	req := Request{
		Path:    "/path",
		Method:  "GET",
		Handler: []gin.HandlerFunc{},
	}

	err := req.validate()
	if err == nil {
		t.Fatal("Request with no handlers should produce error")
	}
	if err.Error() != fmt.Sprintf("request for path %s has no handlers", "/path") {
		t.Fatalf("Expected 'no handlers' error, got: %v", err)
	}
}

func TestRequest_Validate_NilHandlers(t *testing.T) {
	req := Request{
		Path:    "/path",
		Method:  "GET",
		Handler: nil,
	}

	err := req.validate()
	if err == nil {
		t.Fatal("Request with nil handlers should produce error")
	}
}

func TestRequest_Validate_MultipleHandlers(t *testing.T) {
	handler1 := func(c *gin.Context) {}
	handler2 := func(c *gin.Context) {}
	req := Request{
		Path:    "/path",
		Method:  "POST",
		Handler: []gin.HandlerFunc{handler1, handler2},
	}

	err := req.validate()
	if err != nil {
		t.Fatalf("Request with multiple handlers should be valid: %v", err)
	}
}

func TestRequestList_Validate_AllValid(t *testing.T) {
	handler := func(c *gin.Context) {}
	rl := RequestList{
		{Path: "/path1", Method: "GET", Handler: []gin.HandlerFunc{handler}},
		{Path: "/path2", Method: "POST", Handler: []gin.HandlerFunc{handler}},
	}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Valid RequestList should not panic: %v", r)
		}
	}()

	rl.validate()
}

func TestRequestList_Validate_OneInvalid(t *testing.T) {
	handler := func(c *gin.Context) {}
	rl := RequestList{
		{Path: "/path1", Method: "GET", Handler: []gin.HandlerFunc{handler}},
		{Path: "", Method: "POST", Handler: []gin.HandlerFunc{handler}},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("RequestList with invalid request should panic")
		}
	}()

	rl.validate()
}

func TestRequestList_Validate_EmptyList(t *testing.T) {
	rl := RequestList{}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Empty RequestList should not panic: %v", r)
		}
	}()

	rl.validate()
}

func TestPath_AddRequest(t *testing.T) {
	rl := RequestList{}
	path := Path("/test")

	path.addRequest("GET", &rl)

	if len(rl) != 1 {
		t.Fatalf("Expected 1 request after addRequest, got %d", len(rl))
	}
	if rl[0].Path != "/test" {
		t.Fatalf("Expected path '/test', got '%s'", rl[0].Path)
	}
	if rl[0].Method != "GET" {
		t.Fatalf("Expected method 'GET', got '%s'", rl[0].Method)
	}
	if len(rl[0].Handler) != 0 {
		t.Fatalf("Expected 0 handlers initially, got %d", len(rl[0].Handler))
	}
}

func TestPath_AddRequest_MultiplePaths(t *testing.T) {
	rl := RequestList{}
	path1 := Path("/first")
	path2 := Path("/second")

	path1.addRequest("POST", &rl)
	path2.addRequest("PUT", &rl)

	if len(rl) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(rl))
	}
	if rl[0].Path != "/first" {
		t.Fatal("First path mismatch")
	}
	if rl[1].Path != "/second" {
		t.Fatal("Second path mismatch")
	}
}

func TestGinHandler_AddRequest(t *testing.T) {
	handler := func(c *gin.Context) {}
	ginHandler := GinHandler(handler)
	rl := RequestList{
		{Path: "/test", Method: "GET", Handler: make([]gin.HandlerFunc, 0, 1)},
	}

	ginHandler.addRequest("GET", &rl)

	if len(rl[0].Handler) != 1 {
		t.Fatalf("Expected 1 handler, got %d", len(rl[0].Handler))
	}
}

func TestRequestList_AddRoutes_RegistersRoutes(t *testing.T) {
	handler := func(c *gin.Context) {}
	rl := RequestList{
		{Path: "/one", Method: "GET", Handler: []gin.HandlerFunc{handler}},
		{Path: "/two", Method: "POST", Handler: []gin.HandlerFunc{handler}},
		{Path: "/three", Method: "PUT", Handler: []gin.HandlerFunc{handler}},
		{Path: "/four", Method: "DELETE", Handler: []gin.HandlerFunc{handler}},
	}

	router := gin.New()
	rl.AddRoutes(router)

	routes := router.Routes()
	if len(routes) != len(rl) {
		t.Fatalf("Expected %d registered routes, got %d", len(rl), len(routes))
	}

	expected := map[string]struct{}{
		"GET:/one":    {},
		"POST:/two":   {},
		"PUT:/three":  {},
		"DELETE:/four": {},
	}

	for _, route := range routes {
		key := route.Method + ":" + route.Path
		if _, ok := expected[key]; !ok {
			t.Fatalf("Unexpected registered route: %s %s", route.Method, route.Path)
		}
		delete(expected, key)
	}

	if len(expected) != 0 {
		t.Fatalf("Expected routes were not registered: %v", expected)
	}
}

func TestRequestList_AddRoutes_UnsupportedMethodPanics(t *testing.T) {
	handler := func(c *gin.Context) {}
	rl := RequestList{
		{Path: "/bad", Method: "PATCH", Handler: []gin.HandlerFunc{handler}},
	}

	router := gin.New()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for unsupported HTTP method")
		}
	}()

	rl.AddRoutes(router)
}

func TestRequest_MultipleScenarios(t *testing.T) {
	handler := func(c *gin.Context) {}

	tests := []struct {
		name    string
		req     Request
		isValid bool
	}{
		{
			name:    "Valid GET request",
			req:     Request{Path: "/api/users", Method: "GET", Handler: []gin.HandlerFunc{handler}},
			isValid: true,
		},
		{
			name:    "Valid POST request",
			req:     Request{Path: "/api/users", Method: "POST", Handler: []gin.HandlerFunc{handler}},
			isValid: true,
		},
		{
			name:    "Missing path",
			req:     Request{Path: "", Method: "DELETE", Handler: []gin.HandlerFunc{handler}},
			isValid: false,
		},
		{
			name:    "Missing handler",
			req:     Request{Path: "/api", Method: "PUT", Handler: []gin.HandlerFunc{}},
			isValid: false,
		},
		{
			name:    "Parametrized path",
			req:     Request{Path: "/api/users/:id", Method: "GET", Handler: []gin.HandlerFunc{handler}},
			isValid: true,
		},
	}

	for _, test := range tests {
		err := test.req.validate()
		if test.isValid && err != nil {
			t.Errorf("%s: Expected valid, got error: %v", test.name, err)
		}
		if !test.isValid && err == nil {
			t.Errorf("%s: Expected error, got nil", test.name)
		}
	}
}
