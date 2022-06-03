package webapi

import (
	"context"

	"github.com/kordape/tweety/internal/entity"
)

// TweetWebAPI -.
type TweetWebAPI struct {
	accessKey string
	secretKey string
}

// New -.
func New(accessKey string, secretKey string) *TweetWebAPI {
	return &TweetWebAPI{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (t *TweetWebAPI) FetchTweets(ctx context.Context, pageId string) ([]entity.Tweet, error) {
	// TODO: call twitter api to fetch latest tweets from page with pageId
	// parse tweets into internal struct

	tweet1 := entity.Tweet{
		Text: "dummy tweet1",
	}
	tweet2 := entity.Tweet{
		Text: "dummy twee2",
	}
	return []entity.Tweet{tweet1, tweet2}, nil
}
