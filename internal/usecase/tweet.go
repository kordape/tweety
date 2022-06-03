package usecase

import (
	"context"
	"fmt"

	"github.com/kordape/tweety/internal/entity"
)

// TweetUseCase -.
type TweetUseCase struct {
	webAPI TwitterWebAPI
}

// New -.
func New(w TwitterWebAPI) *TweetUseCase {
	return &TweetUseCase{
		webAPI: w,
	}
}

// Classify - classifies if tweets are fake news
func (uc *TweetUseCase) Classify(ctx context.Context) ([]entity.Tweet, error) {
	tweets, err := uc.webAPI.FetchTweets("pageID")
	if err != nil {
		return []entity.Tweet{}, fmt.Errorf("TranslationUseCase - Classify - uc.WebApi.FetchTweets: %w", err)
	}

	return tweets, nil
}
