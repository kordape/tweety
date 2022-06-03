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

}
