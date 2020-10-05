package test

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// AssertMockFullFilled check list of mock to ensure all defined expectations are meet
func AssertMockFullFilled(t *testing.T, mocks ...*mock.Mock) bool {
	res := true
	for _, m := range mocks {
		res = m.AssertExpectations(t) && res
	}
	return res
}

// ResetMocks reset assertions for provided mock list
func ResetMocks(mocks ...*mock.Mock) {
	for _, m := range mocks {
		ResetMock(m)
	}
}

// ResetMock reset assertion for provided mock
func ResetMock(m *mock.Mock) {
	m.ExpectedCalls = []*mock.Call{}
	m.Calls = []mock.Call{}
}
