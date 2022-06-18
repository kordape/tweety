package tweets

import (
	"context"
	"fmt"
	"time"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets/predictor"
	"github.com/kordape/tweety/internal/tweets/twitter"
)

// ensure Classifier implements TweetsClassifier interface
var _ TweetsClassifier = (*Classifier)(nil)

type Classifier struct {
	twitterAPI TwitterWebAPI
	mlModel    MLModel
}

func NewClassifier(t TwitterWebAPI, m MLModel) *Classifier {
	return &Classifier{
		twitterAPI: t,
		mlModel:    m,
	}
}

type ClassifyRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

func (request ClassifyRequest) Validate() error {
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

// Classify - classifies if tweets are fake news
func (classifier *Classifier) Classify(ctx context.Context, request ClassifyRequest) ([]entity.TweetWithClassification, error) {
	tweets, err := classifier.twitterAPI.FetchTweets(ctx, twitter.FetchTweetsRequest{
		MaxResults: request.MaxResults,
		UserId:     request.UserId,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
	})
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("classifier - classify - uc.WebApi.FetchTweets: %w", err)
	}

	predictRequest := make([]predictor.Request, 0)
	for _, t := range tweets {
		predictRequest = append(predictRequest, predictor.Request{Tweet: t.Text})
	}
	predictions, err := classifier.mlModel.FakeTweetPredictor(ctx, predictRequest)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("classifier - classify - uc.WebApi.FetchTweets: %w", err)
	}

	if len(predictions.Prediction) != len(tweets) {
		return []entity.TweetWithClassification{}, fmt.Errorf("not the same number of prediction results (%d) as tweets sent to ml model (%d)", len(predictions.Prediction), len(tweets))
	}

	tweetsWithClassification := []entity.TweetWithClassification{}
	for i, prediction := range predictions.Prediction {
		tweetsWithClassification = append(tweetsWithClassification, entity.TweetWithClassification{
			Text:      tweets[i].Text,
			Fake:      prediction,
			CreatedAt: tweets[i].CreatedAt,
		})
	}

	return tweetsWithClassification, nil
}
