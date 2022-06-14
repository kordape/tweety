package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets"
	"github.com/kordape/tweety/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

const (
	dateLayoutISO = "2006-01-02"
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

func (r *tweetsRoutes) classifyHandler(c *gin.Context) {
	//TODO: add corresponding errorResponse messages
	// refactor naming
	// remove debugging code

	userId, ok := c.Request.URL.Query()["userId"]
	if !ok {
		errorResponse(c, http.StatusBadRequest, "invalid request, invalid userId")
		return
	}
	r.l.Debug("Received userId, %v", userId, ok)

	maxResLimit := 10
	limitQueryParam, ok := c.Request.URL.Query()["maxResults"]
	if ok && len(limitQueryParam) > 0 {
		r.l.Debug("Received maxResults, %v", maxResLimit, ok)
		var err error
		maxResLimit, err = strconv.Atoi(limitQueryParam[0])
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if maxResLimit < 5 || maxResLimit > 100 {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
	}

	var startTime, endTime *time.Time
	fromQueryParams, ok := c.Request.URL.Query()["startTime"]
	if ok && len(fromQueryParams) > 0 {
		r.l.Debug("Received startTime, %v", startTime, ok)
		fromParsed, err := time.Parse(dateLayoutISO, fromQueryParams[0])
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
		startTime = &fromParsed
	}

	toQueryParams, ok := c.Request.URL.Query()["endTime"]
	if ok && len(toQueryParams) > 0 {
		r.l.Debug("Received endTime, %v", endTime, ok)
		toParsed, err := time.Parse(dateLayoutISO, toQueryParams[0])
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
		endTime = &toParsed
	}

	tweets, err := r.t.Classify(c.Request.Context(), userId[0], maxResLimit, startTime, endTime)
	if err != nil {
		r.l.Error(err, "http - v1 - classify")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, classifyResponse{tweets})
}
