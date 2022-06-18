package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets"
	"github.com/kordape/tweety/pkg/logger"
)

const (
	dateLayoutISO     = "2006-01-02"
	defaultMaxResults = 10
)

type tweetsRoutes struct {
	t tweets.TweetsClassifier
	l logger.Interface
}

func newTweetRoutes(handler *gin.RouterGroup, t tweets.TweetsClassifier, l logger.Interface) {
	r := &tweetsRoutes{t, l}

	h := handler.Group("/tweets")
	{
		h.GET("/classify", r.classifyHandler)
	}
}

type classifyResponse struct {
	Tweets []entity.TweetWithClassification `json:"tweets"`
}

func parseClassificationRequest(request *http.Request) (tweets.ClassificationRequest, error) {
	var cr tweets.ClassificationRequest

	userId, ok := request.URL.Query()["userId"]
	if !ok {
		return tweets.ClassificationRequest{}, fmt.Errorf("invalid user id")
	}
	if len(userId) != 1 {
		return tweets.ClassificationRequest{}, fmt.Errorf("invalid user id")
	}
	cr.UserId = userId[0]
	cr.MaxResults = defaultMaxResults
	maxResultsValue, ok := request.URL.Query()["maxResults"]
	if ok && len(maxResultsValue) > 0 {
		var err error
		maxResults, err := strconv.Atoi(maxResultsValue[0])
		if err != nil {
			return tweets.ClassificationRequest{}, fmt.Errorf("error converting string to integer: %w", err)
		}
		cr.MaxResults = maxResults
	}

	startTimeQueryParams, ok := request.URL.Query()["startTime"]
	if ok && len(startTimeQueryParams) > 0 {
		startTimeParsed, err := time.Parse(dateLayoutISO, startTimeQueryParams[0])
		if err != nil {
			return tweets.ClassificationRequest{}, fmt.Errorf("invalid start time parameter")
		}
		cr.StartTime = startTimeParsed.Format(time.RFC3339)
	}

	endTimeQueryParams, ok := request.URL.Query()["endTime"]
	if ok && len(endTimeQueryParams) > 0 {
		endTimeParsed, err := time.Parse(dateLayoutISO, endTimeQueryParams[0])
		if err != nil {
			return tweets.ClassificationRequest{}, fmt.Errorf("invalid end time parameter")
		}
		cr.EndTime = endTimeParsed.Format(time.RFC3339)
	}

	return cr, nil
}

func (r *tweetsRoutes) classifyHandler(c *gin.Context) {
	cr, err := parseClassificationRequest(c.Request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	tweets, err := r.t.Classify(c.Request.Context(), cr)
	if err != nil {
		r.l.Error(err, "")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}
	c.JSON(http.StatusOK, classifyResponse{tweets})
}
