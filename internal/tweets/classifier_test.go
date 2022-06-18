package tweets

import (
	"errors"
	"testing"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func TestClassify(t *testing.T) {
	// tweet := entity.Tweet{
	// 	Text: "tweet",
	// }
	// ctx := context.Background()
	// userId := "1234"
	// twitterApiMock := NewMockTwitterWebAPI(t)
	// twitterApiMock.On("FetchTweets", ctx, userId).Return(
	// 	[]entity.Tweet{
	// 		tweet,
	// 	},
	// 	nil,
	// )
	// classifier := NewClassifier(twitterApiMock)
	// classifiedTweets, err := classifier.Classify(ctx, "1234")
	// assert.NoError(t, err)
	// assert.Equal(t, 1, len(classifiedTweets))
}

func TestClassifyError(t *testing.T) {
	// ctx := context.Background()
	// userId := "1234"
	// twitterApiMock := NewMockTwitterWebAPI(t)
	// twitterApiMock.On("FetchTweets", ctx, userId).Return(
	// 	[]entity.Tweet{},
	// 	errors.New("twitter api error"),
	// )
	// classifier := NewClassifier(twitterApiMock)
	// _, err := classifier.Classify(ctx, "1234")
	// assert.Error(t, err)
}
