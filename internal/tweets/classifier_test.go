package tweets

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets/predictor"
	"github.com/kordape/tweety/internal/tweets/webapi"
)

func TestClassify(t *testing.T) {
	tweet1 := "tweet1"
	tweet2 := "tweet2"

	predictFakeTweetsInput := []predictor.Tweet{
		{Tweet: tweet1}, {Tweet: tweet2},
	}
	ctx := context.Background()
	classificationRequest := ClassificationRequest{
		MaxResults: 10,
		UserId:     "1234",
	}

	fetchTweetsRequest := webapi.FetchTweetsRequest{
		MaxResults: 10,
		UserId:     "1234",
	}

	twitterApiMock := NewMockTwitterWebAPI(t)
	twitterApiMock.On("FetchTweets", ctx, fetchTweetsRequest).Return(
		[]entity.Tweet{
			{Text: tweet1}, {Text: tweet2},
		},
		nil,
	)

	predictorMock := NewMockPredictor(t)
	predictorMock.On("PredictFakeTweets", ctx, predictFakeTweetsInput).Return(
		predictor.Response{Prediction: []int{1, 1}},
		nil,
	)

	classifier := NewClassifier(twitterApiMock, predictorMock)
	classifiedTweets, err := classifier.Classify(ctx, classificationRequest)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(classifiedTweets))
}

func TestClassifyWebAPIError(t *testing.T) {
	ctx := context.Background()
	classificationRequest := ClassificationRequest{
		MaxResults: 10,
		UserId:     "1234",
	}

	fetchTweetsRequest := webapi.FetchTweetsRequest{
		MaxResults: 10,
		UserId:     "1234",
	}
	twitterApiMock := NewMockTwitterWebAPI(t)
	twitterApiMock.On("FetchTweets", ctx, fetchTweetsRequest).Return(
		[]entity.Tweet{},
		errors.New("twitter api error"),
	)
	classifier := NewClassifier(twitterApiMock, nil)
	_, err := classifier.Classify(ctx, classificationRequest)
	assert.Error(t, err)
	assert.EqualError(t, err, "classifier - classify - uc.WebApi.FetchTweets: twitter api error")
}

func TestClassifyPredictorError(t *testing.T) {
	tweet1 := "tweet1"
	tweet2 := "tweet2"

	predictFakeTweetsInput := []predictor.Tweet{
		{Tweet: tweet1}, {Tweet: tweet2},
	}

	ctx := context.Background()
	classificationRequest := ClassificationRequest{
		MaxResults: 10,
		UserId:     "1234",
	}

	fetchTweetsRequest := webapi.FetchTweetsRequest{
		MaxResults: 10,
		UserId:     "1234",
	}
	twitterApiMock := NewMockTwitterWebAPI(t)
	twitterApiMock.On("FetchTweets", ctx, fetchTweetsRequest).Return(
		[]entity.Tweet{
			{Text: tweet1}, {Text: tweet2},
		},
		nil,
	)

	predictorMock := NewMockPredictor(t)
	predictorMock.On("PredictFakeTweets", ctx, predictFakeTweetsInput).Return(
		predictor.Response{},
		errors.New("predictor error"),
	)
	classifier := NewClassifier(twitterApiMock, predictorMock)
	_, err := classifier.Classify(ctx, classificationRequest)
	assert.Error(t, err)
	assert.EqualError(t, err, "classifier - classify - failed to predict fake tweets: predictor error")
}
