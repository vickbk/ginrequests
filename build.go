package ginrequests

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// BuildRequests constructs a slice of Request objects from the provided method and variadic arguments.
// Arguments should alternate between string paths and gin.HandlerFunc handlers.
// Paths (strings) create new Request entries, while handlers are appended to the most recently created request.
//
// Example usage:
//
//	requests := BuildRequests("GET",
//		"/albums", getAlbumsHandler,
//		"/albums/:id", getAlbumByIdHandler,
//		"/health", healthCheckHandler,
//	)
//	// Result: 3 Request objects, each with their respective path and handler
//
// With multiple handlers for the same path:
//
//	requests := BuildRequests("POST",
//		"/albums", validateAlbumHandler, createAlbumHandler,
//		"/health", healthCheckHandler,
//	)
//	// Result: 2 Request objects - first has 2 handlers, second has 1 handler
func BuildRequests(method string, rqs ...any) RequestList {
	requestList := make(RequestList, 0, len(rqs)/2)

	for _, req := range rqs {
		pathOrHandler(req).addRequest(method, &requestList)
	}

	requestList.validate()

	return requestList
}

// isHandlerFunc checks if the given value is a function with signature func(*gin.Context)
func isHandlerFunc(v any) bool {
	if _, ok := v.(gin.HandlerFunc); ok {
		return true
	}
	rt := reflect.TypeOf(v)
	if rt == nil || rt.Kind() != reflect.Func {
		return false
	}

	// Check function signature: must have 1 parameter (*gin.Context) and 0 return values
	if rt.NumIn() != 1 || rt.NumOut() != 0 {
		return false
	}

	// Verify parameter is *gin.Context
	paramType := rt.In(0)
	return paramType.Kind() == reflect.Ptr && paramType.Elem().Name() == "Context" && paramType.Elem().PkgPath() == "github.com/gin-gonic/gin"
}

func pathOrHandler(data any) ReqCaller {
	d, ok := data.(string)
	if ok {
		return Path(d)
	}
	if isHandlerFunc(data) {
		return GinHandler(gin.HandlerFunc(reflect.ValueOf(data).Interface().(func(*gin.Context))))
	}
	panic(fmt.Sprintf(
		"invalid request data: %v - expected string path or gin.HandlerFunc handler.\n"+
			"Example: BuildRequests(\"GET\", \"/path\", handler1, handler2, \"/other\", handler3)",
		data))
}
