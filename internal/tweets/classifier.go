package tweets

import (
	"context"
	"fmt"

	"github.com/kordape/tweety/internal/entity"
)

type Classifier struct {
	webAPI TwitterWebAPI
}

func NewClassfier(w TwitterWebAPI) *Classifier {
	return &Classifier{
		webAPI: w,
	}
}

// Classify - classifies if tweets are fake news
func (classifier *Classifier) Classify(ctx context.Context, userId string) ([]entity.TweetWithClassification, error) {
	tweets, err := classifier.webAPI.FetchTweets(ctx, userId)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("Classifier - Classify - uc.WebApi.FetchTweets: %w", err)
	}

	tweetsWithClassification := []entity.TweetWithClassification{}
	for _, t := range tweets {
		tweetsWithClassification = append(tweetsWithClassification, entity.TweetWithClassification{
			Text: t.Text,
			Fake: 1.0,
		})
	}

	return tweetsWithClassification, nil
}
