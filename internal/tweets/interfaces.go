// Package tweets implements application business logic. Each logic group in own file.
package tweets

import (
	"context"
	v1 "github.com/kordape/tweety/internal/controller/http/v1"
	"github.com/kordape/tweety/internal/entity"
)

// go:generate mockery --name TwitterWebAPI --inpackage --case underscore --filename=./mocks_test.go --disable-version-string
type (
	TweetsClassifier interface {
		Classify(context.Context, v1.ClassifyRequest) ([]entity.TweetWithClassification, error)
	}

	TwitterWebAPI interface {
		FetchTweets(context.Context, v1.ClassifyRequest) ([]entity.Tweet, error)
	}
)
