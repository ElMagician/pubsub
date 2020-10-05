package pubsub_test

import (
	"context"
	"errors"
	"testing"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/elmagician/pubsub"
	mock2 "github.com/elmagician/pubsub/mock"
)

func TestSend(t *testing.T) {
	msgMock := mock2.NewMessages()
	// mockList := []*mock.Mock{&msgMock.Mock}

	ctx := context.Background()

	// init psTest

	Convey("Sending messages to topic", t, func() {
		gcp := pubsub.GCP{
			Config: pubsub.Config{Timeout: 20 * time.Second},
		}

		Convey("should succeed when topics exists and message is correctly formatted", func() {
			srv, cli := initTestServer(ctx)
			defer srv.Close()
			gcp.Client = cli
			So(addTestTopic(ctx, cli, "topic1", "topicAnna", "funkykong"), ShouldBeNil)

			data := []byte("some data for fun")
			attr := map[string]string{"version": "test", "type": "lol"}

			msgMock.On("ToPubsubMessage").Return(attr, data, nil).Once()

			So(gcp.Send(ctx, []string{"topic1", "topicAnna"}, msgMock), ShouldBeNil)
			// So(test.AssertMockFullFilled(&testing.T{}, mockList...), ShouldBeTrue)
		})

		Convey("should errored when message is not marshable", func() {
			srv, cli := initTestServer(ctx)
			defer srv.Close()
			gcp.Client = cli

			So(addTestTopic(ctx, cli, "topic1", "topicAnna", "funkykong"), ShouldBeNil)

			data := []byte("some data for fun")
			attr := map[string]string{"version": "test", "type": "lol"}

			msgMock.
				On("ToPubsubMessage").
				Return(attr, data, errors.New("err")).
				Once()
			err := gcp.Send(ctx, []string{"topic1", "topicAnna"}, msgMock)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, errors.New("err"))
			// So(test.AssertMockFullFilled(&testing.T{}, mockList...), ShouldBeTrue)
		})

		Convey("should errored when topics does not exists", func() {
			srv, cli := initTestServer(ctx)
			defer srv.Close()
			gcp.Client = cli

			data := []byte("some data for fun")
			attr := map[string]string{"version": "test", "type": "lol"}

			msgMock.On("ToPubsubMessage").Return(attr, data, nil).Once()

			So(gcp.Send(ctx, []string{"topics", "topicAnna"}, msgMock), ShouldBeError)
			// So(test.AssertMockFullFilled(&testing.T{}, mockList...), ShouldBeTrue)
		})

		Convey("should errored when send message fail", func() {
			srv, cli := initTestServer(ctx, pstest.WithErrorInjection("Publish", codes.FailedPrecondition, "test ko"))
			defer srv.Close()
			gcp.Client = cli

			So(addTestTopic(ctx, cli, "topic1", "topicAnna", "funkykong"), ShouldBeNil)

			data := []byte("some data for fun")
			attr := map[string]string{"version": "test", "type": "lol"}

			msgMock.On("ToPubsubMessage").Return(attr, data, nil).Once()

			So(gcp.Send(ctx, []string{"topic1", "topicAnna"}, msgMock), ShouldBeError)
			// So(test.AssertMockFullFilled(&testing.T{}, mockList...), ShouldBeTrue)
		})
		//
		// Convey("should errored when message could not be published", func() {
		// 	data := []byte("some data for fun")
		// 	attr := map[string]string{"version": "test", "type": "lol"}
		// 	err := errors.New("could not send message")
		//
		// 	msgMock.On("ToPubsubMessage").Return(attr, data, nil).Once()
		//
		// 	psClient.On("Topic", "topic1").Return(psTopic).Once()
		// 	psTopic.On("Exists", mock.Anything).Return(true, nil).Once()
		//
		// 	psTopic.
		// 		On(
		// 			"Publish",
		// 			mock.Anything,
		// 			pubsubwrapper.AdaptMessage(&googlePubSub.Message{
		// 				Data:       data,
		// 				Attributes: attr,
		// 			}),
		// 		).Return(publishResMock).Once()
		//
		// 	publishResMock.On("Get", mock.Anything).Return("serverID", err).Once()
		//
		// 	So(gcp.Send(ctx, []string{"topic1", "topicAnna"}, msgMock), test.ShouldBeLikeError, pubsub.ErrCouldNotSendMessage)
		// 	So(test.AssertMockFullFilled(&testing.T{}, mockList...), ShouldBeTrue)
		// })

		// test.ResetMocks(mockList...)
	})
}

func initTestServer(ctx context.Context, opts ...pstest.ServerReactorOption) (*pstest.Server, *googlePubSub.Client) {
	srv := pstest.NewServer(opts...)

	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	testClient, err := googlePubSub.NewClient(ctx, "project", option.WithGRPCConn(conn))
	if err != nil {
		panic(err)
	}
	return srv, testClient
}

func addTestTopic(ctx context.Context, cli *googlePubSub.Client, topics ...string) error {
	for _, topic := range topics {
		if _, err := cli.CreateTopic(ctx, topic); err != nil {
			return err
		}

	}
	return nil
}
