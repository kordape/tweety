package tweets

import (
	"context"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
	ml_model "github.com/kordape/tweety/internal/tweets/ml-model"
	"github.com/kordape/tweety/internal/tweets/webapi"
)

type Classifier struct {
	webAPI  TwitterWebAPI
	mlModel MLModel
}

func NewClassifier(w TwitterWebAPI, m MLModel) *Classifier {
	return &Classifier{
		webAPI:  w,
		mlModel: m,
	}
}

// Classify - classifies if tweets are fake news
func (classifier *Classifier) Classify(ctx context.Context, ftr webapi.FetchTweetsRequest) ([]entity.TweetWithClassification, error) {
	tweets, err := classifier.webAPI.FetchTweets(ctx, ftr)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("classifier - classify - uc.WebApi.FetchTweets: %w", err)
	}

	request := make([]ml_model.Tweet, 0)
	for _, t := range tweets {
		request = append(request, ml_model.Tweet{Tweet: t.Text})
	}
	predictions, err := classifier.mlModel.FakeTweetPredictor(ctx, request)
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
