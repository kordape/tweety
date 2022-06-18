// Package tweets implements application business logic. Each logic group in own file.
package tweets

import (
	"context"
	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets/predictor"
	"github.com/kordape/tweety/internal/tweets/webapi"
)

// go:generate mockery --name TwitterWebAPI --inpackage --case underscore --filename=./mocks_test.go --disable-version-string
type (
	TweetsClassifier interface {
		Classify(context.Context, webapi.FetchTweetsRequest) ([]entity.TweetWithClassification, error)
	}

	TwitterWebAPI interface {
		FetchTweets(context.Context, webapi.FetchTweetsRequest) ([]entity.Tweet, error)
	}

	Predictor interface {
		PredictAuthenticTweets(ctx context.Context, tweets []predictor.Tweet) (predictor.Response, error)
	}
)
