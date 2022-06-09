package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
	"io/ioutil"
	"net/http"
)

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

type TwitterWebAPI struct {
	accessKey   string
	secretKey   string
	bearerToken string
}

func New(accessKey string, secretKey string, bearerToken string) *TwitterWebAPI {
	return &TwitterWebAPI{
		accessKey:   accessKey,
		secretKey:   secretKey,
		bearerToken: bearerToken,
	}
}

func (t *TwitterWebAPI) FetchTweets(ctx context.Context, userId string) ([]entity.Tweet, error) {
	// TODO: call twitter api to fetch latest tweets from page with userId
	// http get request to twitter api to fetch tweets
	// parse response into an array of tweet structs
	httpClient := http.Client{}
	// max_results query param for setting the number of tweets to be returned:
	// https://www.postman.com/twitter/workspace/twitter-s-public-workspace/request/9956214-83da6843-c971-4d26-b2f8-07e922a7e285
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(getUsersTweetsUrl, userId), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.bearerToken))
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tweeterResponse getUserTweetsResponse
	err = json.Unmarshal(response, &tweeterResponse)
	if err != nil {
		return nil, err
	}

	return tweeterResponse.Data, nil
}

type getUserTweetsResponse struct {
	Data []entity.Tweet `json:"data"`
	Meta struct {
		ResultCount   int    `json:"result_count"`
		NextToken     string `json:"next_token"`
		PreviousToken string `json:"previous_token"`
	} `json:"meta"`
}
