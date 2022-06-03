// Package tweets implements application business logic. Each logic group in own file.
package tweets

import (
	"context"

	"github.com/kordape/tweety/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	TweetsClassifier interface {
		Classify(context.Context, string) ([]entity.TweetWithClassification, error)
	}

	TwitterWebAPI interface {
		FetchTweets(context.Context, string) ([]entity.Tweet, error)
	}
)
