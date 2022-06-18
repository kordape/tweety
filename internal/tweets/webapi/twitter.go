package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// testing path
// http://localhost:8080/v1/tweets/classify?userId=1277254376&maxResults=90&startTime=2022-01-12&endTime=2022-06-15

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

type TwitterWebAPI struct {
	httpClient  *http.Client
	bearerToken string
}

func New(bearerToken string) *TwitterWebAPI {
	return &TwitterWebAPI{
		bearerToken: bearerToken,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

type FetchTweetsRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

func (ftr FetchTweetsRequest) Validate() error {
	if ftr.MaxResults < 5 || ftr.MaxResults > 100 {
		return fmt.Errorf("invalid max results parameter - can range from 5 to 100")
	}

	if ftr.StartTime != "" && ftr.EndTime != "" {
		start, err := time.Parse(time.RFC3339, ftr.StartTime)
		if err != nil {
			return fmt.Errorf("error parsing start time: %s", err)
		}

		end, err := time.Parse(time.RFC3339, ftr.EndTime)
		if err != nil {
			return fmt.Errorf("error parsing end time: %s", err)
		}

		if start.After(end) {
			return fmt.Errorf("start time is after end time")
		}
	}

	return nil
}

func (t *TwitterWebAPI) FetchTweets(ctx context.Context, ftr FetchTweetsRequest) ([]entity.Tweet, error) {
	baseUrl := fmt.Sprintf(getUsersTweetsUrl, ftr.UserId)
	var queryParams []string
	queryParams = append(queryParams, fmt.Sprintf("max_results=%d", ftr.MaxResults))
	queryParams = append(queryParams, "tweet.fields=id,text,created_at")
	if ftr.StartTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("start_time=%s", ftr.StartTime))
	}
	if ftr.EndTime != "" {
		queryParams = append(queryParams, fmt.Sprintf("end_time=%s", ftr.EndTime))
	}

	url := fmt.Sprintf("%s?%s", baseUrl, strings.Join(queryParams, "&"))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.bearerToken))
	resp, err := t.httpClient.Do(request)
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
