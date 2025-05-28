package functionaltests_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	// TODO have these come from config
	scheme      = "http"
	defaultHost = "localhost"
	port        = 8080
)

type createPostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type postResponseBody struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type createPostResponseBody struct {
	Post postResponseBody `json:"post"`
}

type RestApiTestSuite struct {
	suite.Suite
	baseUrl string
	host    string
	port    int
}

func (suite *RestApiTestSuite) SetupTest() {
	host, ok := os.LookupEnv("REST_API_HOST")
	if !ok {
		host = defaultHost
	}
	suite.baseUrl = fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

func TestRestApiTestSuite(t *testing.T) {
	suite.Run(t, new(RestApiTestSuite))
}
