package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kordape/tweety/internal/entity"
	"github.com/kordape/tweety/internal/usecase"
	"github.com/kordape/tweety/pkg/logger"
)

type tweetRoutes struct {
	t usecase.Tweet
	l logger.Interface
}

func newTweetRoutes(handler *gin.RouterGroup, t usecase.Tweet, l logger.Interface) {
	r := &tweetRoutes{t, l}

	h := handler.Group("/tweet")
	{
		h.GET("/classify", r.classify)
	}
}

type classifyResponse struct {
	Tweets []entity.Tweet `json:"tweets"`
}

func (r *tweetRoutes) classify(c *gin.Context) {
	tweets, err := r.t.Classify(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - classify")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, classifyResponse{tweets})
}
