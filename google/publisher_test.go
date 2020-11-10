package google_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/elmagician/pubsub/google"
)

func ExamplePublisher_Send() {
	var conn *grpc.ClientConn
	var yourEnvelop *Envelop // we are using a mock implementation here
	// This setup a test instance for Google Pubsub.
	// You don't need it outside of unit testing.
	{
		srv, cli := initTestClient("aSuperCoolProject")
		defer func() {
			if err := srv.Close(); err != nil {
				panic(err)
			}
		}()

		_, err := cli.CreateTopic(context.Background(), "mine")
		if err != nil {
			panic(err)
		}
		_, err = cli.CreateTopic(context.Background(), "tropical")
		if err != nil {
			panic(err)
		}
		_, err = cli.CreateTopic(context.Background(), "someTopic")
		if err != nil {
			panic(err)
		}

		conn, err = grpc.Dial(srv.Addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		yourEnvelop = &Envelop{}
		yourEnvelop.message = &message{
			Message: &googlePubSub.Message{
				ID:         "test",
				Data:       []byte("someData"),
				Attributes: map[string]string{"version": "v1", "type": "test", "new": "false"},
			},
		}
	}

	ctx := context.Background()
	conf := google.Config{
		ProjectID:       "aSuperCoolProject",
		CredentialsPath: "path/to/credentials.json",
		Timeout:         10 * time.Second,
		Concurrency:     0,
	}

	// Initialize pubsub instance.
	ps, err := google.NewPubsub(ctx, conf, option.WithGRPCConn(conn))
	if err != nil {
		// TODO manage error
		panic(err)
		return
	}

	// Register topics that will be used from the instance.
	ps.Registry().
		MustAddTopic("topicAnna", nil). // Will not fail if topic does not exists using MustAddTopic. Use AddTopic to check existence on add.
		MustAddTopic("tropical", nil)

	if err := ps.Registry().AddTopic("mine", nil); err != nil {
		// TODO manager error
		panic(err)
		return
	}

	// Send message to registered topics
	results, err := ps.Publish().To("mine").Send(ctx, yourEnvelop)
	if err != nil {
		// TODO manage error
		fmt.Println(err.Error())
		return
	}

	idsStr := ""
	for topic, res := range results.Results(context.Background()) {
		idsStr += topic + ": " + res.ID
	}

	idsStr = strings.TrimSuffix(idsStr, ",")
	fmt.Println("Msg:", idsStr)

	// Send message to registered topics
	results, err = ps.Publish().To("tropical", "mine").Send(ctx, yourEnvelop)
	if err != nil {
		// TODO manage error
		fmt.Println(err.Error())
		return
	}

	idsStr = ""
	for topic, res := range results.Results(context.Background()) {
		if res.Error != nil {
			fmt.Println("got unexpected error", res.Error)
		}
		if topic == "tropical" {
			idsStr = topic + ", " + idsStr
		} else {
			idsStr += topic
		}
	}

	idsStr = strings.TrimSuffix(idsStr, ",")
	fmt.Println(idsStr)
	// Stop topics connection to server without discarding them.s
	ps.Registry().StopTopics("mine", "someTopic")

	// Reset registry after stopping all topics.
	ps.Registry().Clear()

	// Output:
	// Msg: mine: m0
	// tropical, mine
}
