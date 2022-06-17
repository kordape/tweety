package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets"
	"io/ioutil"
	"net/http"
	"strings"
)

// testing path
// http://localhost:8080/v1/tweets/classify?userId=1277254376&maxResults=90&startTime=2022-01-12&endTime=2022-06-15

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

type TwitterWebAPI struct {
	bearerToken string
}

func New(bearerToken string) *TwitterWebAPI {
	return &TwitterWebAPI{
		bearerToken: bearerToken,
	}
}

func (t *TwitterWebAPI) FetchTweets(ctx context.Context, classifyRequest tweets.ClassifyRequest) ([]entity.Tweet, error) {
	baseUrl := fmt.Sprintf(getUsersTweetsUrl, classifyRequest.UserId)
	var queryParams []string
	queryParams = append(queryParams, fmt.Sprintf("max_results=%d", classifyRequest.MaxResults))
	queryParams = append(queryParams, "tweet.fields=id,text,created_at")
	if classifyRequest.StartTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("start_time=%s", classifyRequest.StartTime))
	}
	if classifyRequest.EndTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("end_time=%s", classifyRequest.EndTime))
	}

	url := fmt.Sprintf("%s?%s", baseUrl, strings.Join(queryParams, "&"))
	httpClient := http.Client{}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.bearerToken))
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var tweeterResponse getUserTweetsResponse
	err = json.Unmarshal(response, &tweeterResponse)
	if err != nil {
		return nil, err
	}

	return tweeterResponse.Data, nil
}

type getUserTweetsResponse struct {
	Data []entity.Tweet         `json:"data"`
	Meta TweetsResponseMetaData `json:"meta"`
}

// TweetsResponseMetaData left to enable pagination option in perspective
// can be removed if needed
type TweetsResponseMetaData struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}
