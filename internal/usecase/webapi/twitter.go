package webapi

import (
	"github.com/kordape/tweety/internal/entity"
)

// TweetWebAPI -.
type TweetWebAPI struct {
	conf string
}

// New -.
func New() *TweetWebAPI {
	return &TweetWebAPI{
		conf: "conf",
	}
}

// FetchTweets -.
func (t *TweetWebAPI) FetchTweets(pageId string) ([]entity.Tweet, error) {
	return []entity.Tweet{
		entity.Tweet{
			Body: "dummy tweet",
		},
		entity.Tweet{
			Body: "dummy tweet1",
		},
	}, nil
}
