package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alec-w/sample-applications-and-libraries/libraries/logging"
)

// server is internal implmentation of generated strict API interface
type server struct {
	logger logging.Logger
}

// newServer instantiates the internal implementation of generated strict API interface
func newServer(logger logging.Logger) *server {
	return &server{logger: logger}
}

// NewServer returns an *http.Server that fulfils the API spec
func NewServer(port int, logger logging.Logger) *http.Server {
	return &http.Server{
		Handler: HandlerFromMux(NewStrictHandler(newServer(logger), nil), http.NewServeMux()),
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
}

// GetPosts implements (GET /posts) API endpoint
func (s *server) GetPosts(ctx context.Context, request GetPostsRequestObject) (GetPostsResponseObject, error) {
	s.logger.Info("Serving list posts request")
	resp := GetPosts200JSONResponse{
		Posts: []PostResponse{
			{
				Id:        "one",
				Title:     "test",
				Content:   "Test Post",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	return resp, nil
}
