package test_test

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/elmagician/pubsub/test"
)

func TestResetMock(t *testing.T) {
	Convey("Given some mock", t, func() {
		Convey("should be able to reset one", func() {
			initTestStack()
			So(testMock.Test("test"), ShouldBeNil)
			So(testMock.Test("another"), ShouldResemble, errors.New("err2"))

			test.ResetMock(&testMock.Mock)
			testMock.On("Test", "test").Return(errors.New("random"))

			So(testMock.Test("test"), ShouldResemble, errors.New("random"))
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("err"))
		})

		Convey("should be able to reset list", func() {
			initTestStack()
			So(testMock.Test("test"), ShouldBeNil)
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("err"))
			So(testMock.Test("another"), ShouldResemble, errors.New("err2"))

			test.ResetMocks(&testMock.Mock, &testMock2.Mock)
			testMock.On("Test", "test").Return(errors.New("random"))
			testMock2.On("Do", "smtg").Return(errors.New("some"))

			So(testMock.Test("test"), ShouldResemble, errors.New("random"))
			So(testMock2.Do("smtg"), ShouldResemble, errors.New("some"))
		})

		test.ResetMocks(&testMock.Mock, &testMock2.Mock)
	})
}

func TestAssertMockFullFilled(t *testing.T) {
	Convey("Given some mock", t, func() {
		initTestStack()
		tester := &testing.T{}
		Convey("should be ok if all expectation were met", func() {
			_ = testMock.Test("test")
			_ = testMock2.Do("smtg")
			_ = testMock.Test("another")

			So(test.AssertMockFullFilled(tester, &testMock.Mock, &testMock2.Mock), ShouldBeTrue)
			So(tester.Failed(), ShouldBeFalse)
		})

		Convey("should be false if some expectations were not met", func() {
			_ = testMock.Test("test")
			So(test.AssertMockFullFilled(tester, &testMock.Mock), ShouldBeFalse)
			So(test.AssertMockFullFilled(tester, &testMock2.Mock), ShouldBeFalse)
			So(tester.Failed(), ShouldBeTrue)

			_ = testMock2.Do("smtg")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock), ShouldBeTrue)
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeFalse)

			_ = testMock.Test("another")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeTrue)
			So(tester.Failed(), ShouldBeTrue)

			test.ResetMocks(&testMock.Mock, &testMock2.Mock)
			initTestStack()

			_ = testMock.Test("test")
			So(test.AssertMockFullFilled(tester, &testMock2.Mock, &testMock.Mock), ShouldBeFalse)
			_ = testMock.Test("another")
			So(test.AssertMockFullFilled(tester, &testMock.Mock), ShouldBeTrue)
			So(test.AssertMockFullFilled(tester, &testMock.Mock, &testMock2.Mock), ShouldBeFalse)
			So(tester.Failed(), ShouldBeTrue)
		})

		test.ResetMocks(&testMock.Mock, &testMock2.Mock)
	})
}

type Mock struct {
	mock.Mock
}

func (m *Mock) Test(s string) error {
	args := m.Called(s)
	return args.Error(0)
}

type Mock2 struct {
	mock.Mock
}

func (m *Mock2) Do(action string) error {
	args := m.Called(action)
	return args.Error(0)
}

var (
	testMock  = &Mock{}
	testMock2 = &Mock2{}
)

func initTestStack() {
	testMock.On("Test", "test").Return(nil)
	testMock2.On("Do", "smtg").Return(errors.New("err"))
	testMock.On("Test", "another").Return(errors.New("err2"))
}
