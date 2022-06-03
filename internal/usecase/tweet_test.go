package usecase_test

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

// func classify(t *testing.T) (*usecase.TweetUseCase, *MockTranslationWebAPI) {
// 	t.Helper()

// 	mockCtl := gomock.NewController(t)
// 	defer mockCtl.Finish()

// 	webAPI := NewMockTranslationWebAPI(mockCtl)

// 	translation := usecase.New(webAPI)

// 	return translation, webAPI
// }

func TestClassify(t *testing.T) {

}
