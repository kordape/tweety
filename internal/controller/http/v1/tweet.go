package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/tweets"
	"github.com/kordape/tweety/pkg/logger"
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
	pageId, ok := c.Request.URL.Query()["pageId"]
	r.l.Debug("Received pageId", pageId, ok)

	if !ok {
		errorResponse(c, http.StatusBadRequest, "invalid request")

		return
	}

	tweets, err := r.t.Classify(c.Request.Context(), pageId[0])
	if err != nil {
		r.l.Error(err, "http - v1 - classify")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, classifyResponse{tweets})
}
