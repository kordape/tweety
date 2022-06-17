package tweets

import (
	"context"
	"fmt"
	"github.com/kordape/tweety/internal/entity"
)

type Classifier struct {
	webAPI TwitterWebAPI
}

type ClassifyRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

func NewClassfier(w TwitterWebAPI) *Classifier {
	return &Classifier{
		webAPI: w,
	}
}

// Classify - classifies if tweets are fake news
func (classifier *Classifier) Classify(ctx context.Context, cr ClassifyRequest) ([]entity.TweetWithClassification, error) {
	tweets, err := classifier.webAPI.FetchTweets(ctx, cr)
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
