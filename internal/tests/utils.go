package tests

import (
	"fmt"
	"testing"

	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
)

func Assert(t *testing.T, predicate bool, failMessage string) {
	if !predicate {
		logger.Log("Test failed:", failMessage)
		t.FailNow()
	}
}

func AssertMerrorDoesNotExist(t *testing.T, error *merrors.MError) {
	if nil != error {
		logger.Error("Test failed with error:", error.Message)
		t.FailNow()
	}
}

func AssertMerror(t *testing.T, error *merrors.MError, code merrors.MErrorCode, message string) {
	if nil == error {
		logger.Error("Test failed because error is empty.")
		t.FailNow()
	}
	Assert(t, error.Code == code && error.Message == message, fmt.Sprintf("\n[%d - %s] \nwas expected but \n[%d - %s] \nwas found\n", code, message, error.Code, error.Message))
}
