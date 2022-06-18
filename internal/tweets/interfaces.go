// Package tweets implements application business logic. Each logic group in own file.
package tweets

import (
	"context"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets/predictor"
	"github.com/kordape/tweety/internal/tweets/twitter"
)

// go:generate mockery --name TwitterWebAPI --inpackage --case underscore --filename=./mocks_test.go --disable-version-string
type (
	TweetsClassifier interface {
		Classify(context.Context, ClassifyRequest) ([]entity.TweetWithClassification, error)
	}

	TwitterWebAPI interface {
		FetchTweets(context.Context, twitter.FetchTweetsRequest) ([]entity.Tweet, error)
	}

	MLModel interface {
		FakeTweetPredictor(ctx context.Context, tweets []predictor.Request) (predictor.Response, error)
	}
)
