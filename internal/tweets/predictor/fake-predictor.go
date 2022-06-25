package predictor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const predictUrl = "http://ml:8080/predict"

type Predictor struct {
	httpClient *http.Client
}

func New() *Predictor {
	return &Predictor{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type Tweet struct {
	Tweet string `json:"tweet"`
}
type Response struct {
	Prediction []int `json:"prediction"`
}

func (p *Predictor) PredictFakeTweets(ctx context.Context, tweets []Tweet) (Response, error) {
	buf, err := json.Marshal(tweets)
	if err != nil {
		return Response{}, fmt.Errorf("error marshalling request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, predictUrl, bytes.NewBuffer(buf))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return Response{}, fmt.Errorf("error creating http request: %w", err)
	}

	resp, err := p.httpClient.Do(request)
	if err != nil {
		return Response{}, fmt.Errorf("error doing http request: %w", err)
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading response: %w", err)
	}

	var fpr Response
	err = json.Unmarshal(response, &fpr)
	if err != nil {
		return Response{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return fpr, nil
}
