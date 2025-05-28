package functionaltests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (suite *RestApiTestSuite) createPost(ctx context.Context, input createPostRequestBody) (int, []byte) {
	requestBody, err := json.Marshal(input)
	suite.Require().NoError(err)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/posts", suite.baseUrl),
		bytes.NewReader(requestBody),
	)
	suite.Require().NoError(err)

	response, err := http.DefaultClient.Do(request)
	suite.Require().NoError(err)

	responseBody, err := io.ReadAll(response.Body)
	suite.Require().NoError(err)
	defer response.Body.Close()

	return response.StatusCode, responseBody
}

func (suite *RestApiTestSuite) TestCreateSuccess() {
	beforeRequest := time.Now()
	input := createPostRequestBody{Title: "test title", Content: "some content to see it worked"}
	statusCode, responseBody := suite.createPost(suite.T().Context(), input)
	afterRequest := time.Now()

	if !suite.Assert().Equal(http.StatusCreated, statusCode) {
		suite.T().Log(string(responseBody))
		suite.T().FailNow()
	}

	var output createPostResponseBody
	if !suite.Assert().NoError(json.Unmarshal(responseBody, &output)) {
		suite.T().Log(string(responseBody))
		suite.T().FailNow()
	}

	suite.Assert().GreaterOrEqual(output.Post.Id, 1)
	suite.Assert().Equal(input.Title, output.Post.Title)
	suite.Assert().Equal(input.Content, output.Post.Content)
	suite.Assert().GreaterOrEqual(output.Post.CreatedAt.Unix(), beforeRequest.Unix())
	suite.Assert().LessOrEqual(output.Post.CreatedAt.Unix(), afterRequest.Unix())
}
