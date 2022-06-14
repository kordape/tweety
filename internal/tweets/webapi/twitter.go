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

func (t *TwitterWebAPI) FetchTweets(ctx context.Context, userId string, maxResults int, from *time.Time, to *time.Time) ([]entity.Tweet, error) {
	// http get request to twitter api to fetch tweets
	// parse response into an array of tweet structs
	baseUrl := fmt.Sprintf(getUsersTweetsUrl, userId)
	var queryParams []string
	// max_results query param for setting the number of tweets to be returned: min=5 max =100
	queryParams = append(queryParams, fmt.Sprintf("max_results=%d", maxResults))
	queryParams = append(queryParams, "tweet.fields=id,text,created_at")
	if from != nil {
		tmp := *from
		queryParams = append(queryParams, fmt.Sprintf("start_time=%s", tmp.Format(time.RFC3339)))
	}
	if to != nil {
		tmp := *to
		queryParams = append(queryParams, fmt.Sprintf("end_time=%s", tmp.Format(time.RFC3339)))
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
	Data []entity.Tweet `json:"data"`
	Meta struct {
		ResultCount   int    `json:"result_count"`
		NextToken     string `json:"next_token"`
		PreviousToken string `json:"previous_token"`
	} `json:"meta"`
}
