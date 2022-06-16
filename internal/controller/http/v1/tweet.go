package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets"
	"github.com/kordape/tweety/pkg/logger"
	"net/http"
	"strconv"
	"time"
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

type ClassifyRequest struct {
	MaxResults int
	UserId     string
	StartTime  string
	EndTime    string
}

func parseRequest(request *http.Request) (*ClassifyRequest, error) {
	var classifyRequest ClassifyRequest

	userId, ok := request.URL.Query()["userId"]
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	if len(userId) == 1 {
		classifyRequest.UserId = userId[0]
	}

	maxResults := defaultMaxResults
	maxResultsValue, ok := request.URL.Query()["maxResults"]
	if ok && len(maxResultsValue) > 0 {
		var err error
		maxResults, err = strconv.Atoi(maxResultsValue[0])
		if err != nil {
			return nil, fmt.Errorf("error converting string to integer: %w", err)
		}
		if maxResults < 5 || maxResults > 100 {
			return nil, fmt.Errorf("invalid max results parameter - can range from 5 to 100")
		}
	}
	classifyRequest.MaxResults = maxResults

	startTimeQueryParams, ok := request.URL.Query()["startTime"]
	if ok && len(startTimeQueryParams) > 0 {
		startTimeParsed, err := time.Parse(dateLayoutISO, startTimeQueryParams[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start time parameter")
		}
		classifyRequest.StartTime = startTimeParsed.Format(time.RFC3339)
	}

	endTimeQueryParams, ok := request.URL.Query()["endTime"]
	if ok && len(endTimeQueryParams) > 0 {
		endTimeParsed, err := time.Parse(dateLayoutISO, endTimeQueryParams[0])
		if err != nil {
			return nil, fmt.Errorf("invalid end time parameter")
		}
		classifyRequest.EndTime = endTimeParsed.Format(time.RFC3339)
	}
	if classifyRequest.StartTime > classifyRequest.EndTime {
		return nil, fmt.Errorf("invalid time parameters - start time after end time")
	}
	return &classifyRequest, nil
}

func (r *tweetsRoutes) classifyHandler(c *gin.Context) {
	cr, err := parseRequest(c.Request)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	tweets, err := r.t.Classify(c.Request.Context(), *cr)
	if err != nil {
		r.l.Error(err, "http - v1 - classify")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}
	c.JSON(http.StatusOK, classifyResponse{tweets})
}
