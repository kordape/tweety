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
	userId, ok := c.Request.URL.Query()["userId"]
	if !ok {
		errorResponse(c, http.StatusBadRequest, "invalid request, userId is missing")
		return
	}
	r.l.Debug("Received userId, %v", userId, ok)

	limit := 10
	var from, to *time.Time

	limitQueryParam, ok := c.Request.URL.Query()["limit"]
	if ok && len(limitQueryParam) > 0 {
		r.l.Debug("Received numberOfResults, %v", limit, ok)
		var err error
		limit, err = strconv.Atoi(limitQueryParam[0])
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if limit < 5 || limit > 100 {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
	}

	fromQueryParams, ok := c.Request.URL.Query()["from"]
	if ok && len(fromQueryParams) > 0 {
		r.l.Debug("Received from, %v", from, ok)
		fromParsed, err := time.Parse(dateLayoutISO, fromQueryParams[0])
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
		from = &fromParsed
	}

	toQueryParams, ok := c.Request.URL.Query()["to"]
	if ok && len(toQueryParams) > 0 {
		r.l.Debug("Received to, %v", to, ok)
		toParsed, err := time.Parse(dateLayoutISO, toQueryParams[0])
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}
		to = &toParsed
	}

	tweets, err := r.t.Classify(c.Request.Context(), userId[0], limit, from, to)
	if err != nil {
		r.l.Error(err, "http - v1 - classify")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, classifyResponse{tweets})
}
