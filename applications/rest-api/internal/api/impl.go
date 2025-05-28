package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alec-w/sample-applications-and-libraries/applications/rest-api/internal/models"
	"github.com/alec-w/sample-applications-and-libraries/libraries/logging"
)

type Datastore interface {
	CreatePost(ctx context.Context, title, content string, createdAt time.Time) (models.Post, error)
	ListPosts(ctx context.Context) ([]models.Post, error)
}

// server is internal implmentation of generated strict API interface
type server struct {
	datastore Datastore
	logger    logging.Logger
}

// newServer instantiates the internal implementation of generated strict API interface
func newServer(datastore Datastore, logger logging.Logger) *server {
	return &server{datastore: datastore, logger: logger}
}

// NewServer returns an *http.Server that fulfils the API spec
func NewServer(port int, logger logging.Logger, datastore Datastore) *http.Server {
	return &http.Server{
		Handler: HandlerFromMux(NewStrictHandler(newServer(datastore, logger), nil), http.NewServeMux()),
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
}

// GetPosts implements (GET /posts) API endpoint
func (s *server) GetPosts(ctx context.Context, request GetPostsRequestObject) (GetPostsResponseObject, error) {
	s.logger.Debug("Serving list posts request")
	posts, err := s.datastore.ListPosts(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to list posts")
		return GetPosts500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Message: "An internal error occurred.",
			},
		}, nil
	}
	resp := GetPosts200JSONResponse{
		Posts: postsToApiResponses(posts),
	}

	return resp, nil
}

// PostPosts implements (POST /posts) API endpoint
func (s *server) PostPosts(ctx context.Context, request PostPostsRequestObject) (PostPostsResponseObject, error) {
	s.logger.Debug("Serving create post request")
	post, err := s.datastore.CreatePost(ctx, request.Body.Title, request.Body.Content, time.Now())
	if err != nil {
		s.logger.WithError(err).Error("failed to list posts")
		return PostPosts500JSONResponse{
			InternalErrorJSONResponse: InternalErrorJSONResponse{
				Message: "An internal error occurred.",
			},
		}, nil
	}
	resp := PostPosts201JSONResponse{
		Post: postToApiResponse(post),
	}

	return resp, nil
}

func postToApiResponse(internalModel models.Post) PostResponse {
	return PostResponse{
		Id:        internalModel.Id,
		Title:     internalModel.Title,
		Content:   internalModel.Content,
		CreatedAt: internalModel.CreatedAt.Format(time.RFC3339),
	}
}

func postsToApiResponses(internalModels []models.Post) []PostResponse {
	responses := make([]PostResponse, len(internalModels))
	for i, post := range internalModels {
		responses[i] = postToApiResponse(post)
	}
	return responses
}
