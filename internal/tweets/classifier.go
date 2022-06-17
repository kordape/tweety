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

type ClassifyRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

func (cr ClassifyRequest) Validate() error {
	if cr.MaxResults < 5 || cr.MaxResults > 100 {
		return fmt.Errorf("invalid max results parameter - can range from 5 to 100")
	}

	if cr.StartTime != "" && cr.EndTime != "" {
		start, err := time.Parse(time.RFC3339, cr.StartTime)
		if err != nil {
			return fmt.Errorf("error parsing start time: %s", err)
		}

		end, err := time.Parse(time.RFC3339, cr.EndTime)
		if err != nil {
			return fmt.Errorf("error parsing end time: %s", err)
		}

		if start.After(end) {
			return fmt.Errorf("start time is after end time")
		}
	}

	return nil
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
