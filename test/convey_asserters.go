package test

import (
	"fmt"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/mock"
)

// Assert errors looks alike using goerrors.Is method
func ShouldBeLikeError(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return "ShouldBeLikeError assertion expect 1 and only 1 asserter"
	}

	resp := ""
	expActComp := fmt.Sprintf("\nActual: %#v\nExpected: %#v", actual, expected[0])

	actualErr, ok := actual.(error)
	if !ok {
		resp += "Actual value should be an error"
	}

	expectedErr, ok2 := expected[0].(error)
	if !ok2 {
		resp += "Expected value should be an error"
	}

	if resp != "" {
		return resp + expActComp
	}

	if !errors.Is(actualErr, expectedErr) {
		return "Actual error do not match expected err Is." + expActComp
	}

	return ""
}

// Assert provided DB/testify mock is full filled
func ShouldBeFullFilled(actual interface{}, _ ...interface{}) string {
	expActComp := fmt.Sprintf("\nProvided: %#v", actual)

	singleTestifyMock, isTestify := actual.(*mock.Mock)
	listTestifyMock, isTestifyList := actual.([]*mock.Mock)

	if !isTestify && !isTestifyList {
		return "Provided value is not a testify mock or mock list." + expActComp
	}

	if isTestify && !AssertMockFullFilled(&testing.T{}, singleTestifyMock) {
		return "Some expectation were not met."
	}

	if isTestifyList {
		for pos, localMock := range listTestifyMock {
			if !AssertMockFullFilled(&testing.T{}, localMock) {
				return fmt.Sprintf("Some expectation were not met for mock number: %d.", pos)
			}
		}
	}

	return ""
}
