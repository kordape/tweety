package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kordape/tweety/internal/entity"
)

// testing path
// http://localhost:8080/v1/tweets/classify?userId=1277254376&maxResults=90&startTime=2022-01-12&endTime=2022-06-15

const (
	getUsersTweetsUrl = "https://api.twitter.com/2/users/%s/tweets/"
)

type Client struct {
	httpClient  *http.Client
	bearerToken string
}

func New(bearerToken string) *Client {
	return &Client{
		bearerToken: bearerToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type FetchTweetsRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

type getUserTweetsResponse struct {
	Data []Tweet  `json:"data"`
	Meta Metadata `json:"meta"`
}

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Text      string `json:"text"`
}

// metaData left to enable pagination option in perspective
// can be removed if needed
type Metadata struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

func (request FetchTweetsRequest) Validate() error {
	if request.MaxResults < 5 || request.MaxResults > 100 {
		return fmt.Errorf("invalid max results parameter - can range from 5 to 100")
	}

	if request.StartTime != "" && request.EndTime != "" {
		start, err := time.Parse(time.RFC3339, request.StartTime)
		if err != nil {
			return fmt.Errorf("error parsing start time: %s", err)
		}

		end, err := time.Parse(time.RFC3339, request.EndTime)
		if err != nil {
			return fmt.Errorf("error parsing end time: %s", err)
		}

		if start.After(end) {
			return fmt.Errorf("start time is after end time")
		}
	}

	return nil
}

func (client *Client) FetchTweets(ctx context.Context, ftr FetchTweetsRequest) ([]entity.Tweet, error) {
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
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.bearerToken))
	resp, err := client.httpClient.Do(request)
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

	result := make([]entity.Tweet, len(tweeterResponse.Data))
	for i, tweet := range tweeterResponse.Data {
		result[i] = entity.Tweet{
			Id:        tweet.Id,
			Text:      tweet.Text,
			CreatedAt: tweet.CreatedAt,
		}
	}

	return result, nil
}
