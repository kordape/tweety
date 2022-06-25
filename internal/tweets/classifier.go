package tweets

import (
	"context"
	"fmt"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets/predictor"
	"github.com/kordape/tweety/internal/tweets/webapi"
)

type Classifier struct {
	webAPI    TwitterWebAPI
	predictor Predictor
}

func NewClassifier(w TwitterWebAPI, p Predictor) *Classifier {
	return &Classifier{
		webAPI:    w,
		predictor: p,
	}
}

// Classify - classifies if tweets are fake news
func (classifier *Classifier) Classify(ctx context.Context, cr ClassificationRequest) ([]entity.TweetWithClassification, error) {
	ftr := webapi.FetchTweetsRequest{
		MaxResults: cr.MaxResults,
		UserId:     cr.UserId,
		StartTime:  cr.StartTime,
		EndTime:    cr.EndTime,
	}

	err := ftr.Validate()
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("validation of FetchTweetsRequest failed: %w", err)
	}

	tweets, err := classifier.webAPI.FetchTweets(ctx, ftr)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("classifier - classify - uc.WebApi.FetchTweets: %w", err)
	}

	request := make([]predictor.Tweet, 0)
	for _, t := range tweets {
		request = append(request, predictor.Tweet{Tweet: t.Text})
	}
	predictions, err := classifier.predictor.PredictFakeTweets(ctx, request)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("classifier - classify - failed to predict fake tweets: %w", err)
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

type ClassificationRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}
