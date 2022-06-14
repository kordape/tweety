package tweets

import (
	"context"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
	"time"
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
func (classifier *Classifier) Classify(ctx context.Context, userId string, maxResults int, from, to *time.Time) ([]entity.TweetWithClassification, error) {
	tweets, err := classifier.webAPI.FetchTweets(ctx, userId, maxResults, from, to)
	if err != nil {
		return []entity.TweetWithClassification{}, fmt.Errorf("Classifier - Classify - uc.WebApi.FetchTweets: %w", err)
	}

	tweetsWithClassification := []entity.TweetWithClassification{}
	for _, t := range tweets {
		tweetsWithClassification = append(tweetsWithClassification, entity.TweetWithClassification{
			Text:      t.Text,
			Fake:      1.0,
			CreatedAt: t.CreatedAt,
		})
	}

	return tweetsWithClassification, nil
}
