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

func parseRequest(request *http.Request, r *tweetsRoutes) (*ClassifyRequest, error) {
	var classifyRequest ClassifyRequest
	userId, ok := request.URL.Query()["userId"]
	if !ok {
		return nil, fmt.Errorf("can't find user")
	}
	//TODO if len !=1
	classifyRequest.UserId = userId[0]
	maxResults := defaultMaxResults
	maxResultsValue, ok := request.URL.Query()["maxResults"]
	if ok && len(maxResultsValue) > 0 {
		var err error
		maxResults, err = strconv.Atoi(maxResultsValue[0])
		if err != nil {
			return nil, fmt.Errorf("error converting string to integer: %w", err)
		}
		if maxResults < 5 || maxResults > 100 {
			//TODO add proper error
			return nil, fmt.Errorf("...")
		}
	}
	classifyRequest.MaxResults = maxResults
	fromQueryParams, ok := request.URL.Query()["startTime"]
	if ok && len(fromQueryParams) > 0 {
		fromParsed, err := time.Parse(dateLayoutISO, fromQueryParams[0])
		if err != nil {
			//TODO add proper error
			return nil, fmt.Errorf("...")
		}

		classifyRequest.StartTime = fromParsed.Format(time.RFC3339)
	}

	toQueryParams, ok := request.URL.Query()["endTime"]
	if ok && len(toQueryParams) > 0 {
		toParsed, err := time.Parse(dateLayoutISO, toQueryParams[0])
		if err != nil {
			//TODO add proper error
			return nil, fmt.Errorf("...")
		}
		classifyRequest.EndTime = toParsed.Format(time.RFC3339)
	}
	return &classifyRequest, nil
}

func (r *tweetsRoutes) classifyHandler(c *gin.Context) {
	//TODO: add corresponding errorResponse messages
	// refactor naming
	// remove debugging code

	cr, err := parseRequest(c.Request, r)
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
