package ginrequests

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Handler      []gin.HandlerFunc
	Method, Path string
}

// validate checks that the Request has required fields
func (req *Request) validate() error {
	if req.Path == "" {
		return fmt.Errorf("request path cannot be empty")
	}
	if len(req.Handler) == 0 {
		return fmt.Errorf("request for path %s has no handlers", req.Path)
	}
	return nil
}

type RequestList []Request

func (rl *RequestList) validate() {
	for _, req := range *rl {
		if err := req.validate(); err != nil {
			panic(fmt.Sprintf("invalid request configuration: %v", err))
		}
	}
}

// AddRoutes registers each request in the list with the provided Gin router.
//
// Supported methods are GET, POST, PUT, and DELETE. If a request contains an
// unsupported HTTP method, AddRoutes will panic.
func (rl *RequestList) AddRoutes(router *gin.Engine) {
	for _, req := range *rl {
		switch req.Method {
		case "GET":
			router.GET(req.Path, req.Handler...)
		case "POST":
			router.POST(req.Path, req.Handler...)
		case "PUT":
			router.PUT(req.Path, req.Handler...)
		case "DELETE":
			router.DELETE(req.Path, req.Handler...)
		default:
			panic(fmt.Sprintf("unsupported HTTP method: %s", req.Method))
		}
	}
}

type ReqCaller interface {
	addRequest(method string, requests *RequestList)
}

type Path string

func (p Path) addRequest(method string, requests *RequestList) {
	*requests = append(*requests, Request{
		Handler: make([]gin.HandlerFunc, 0, 1),
		Method:  method,
		Path:    string(p),
	})
}

type GinHandler gin.HandlerFunc

func (h GinHandler) addRequest(_ string, requests *RequestList) {
	if len(*requests) == 0 {
		panic("handler added before any path - paths must precede handlers")
	}
	lastIdx := len(*requests) - 1
	(*requests)[lastIdx].Handler = append((*requests)[lastIdx].Handler, gin.HandlerFunc(h))
}
