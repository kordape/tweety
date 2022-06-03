// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/kordape/tweety/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Tweet -.
	Tweet interface {
		Classify(context.Context) ([]entity.Tweet, error)
	}

	// TranslationWebAPI -.
	TwitterWebAPI interface {
		FetchTweets(string) ([]entity.Tweet, error)
	}
)
