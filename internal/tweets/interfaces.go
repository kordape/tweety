// Package tweets implements application business logic. Each logic group in own file.
package tweets

import (
	"context"

	"github.com/kordape/tweety/internal/entity"
)

// go:generate mockery --name TwitterWebAPI --inpackage --case underscore --filename=./mocks_test.go --disable-version-string
type (
	TweetsClassifier interface {
		Classify(context.Context, string) ([]entity.TweetWithClassification, error)
	}

	TwitterWebAPI interface {
		FetchTweets(context.Context, string) ([]entity.Tweet, error)
	}
)
