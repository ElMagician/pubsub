package google_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/api/option"

	"github.com/elmagician/pubsub/google"
)

func TestPubsub_Publisher(t *testing.T) {
	Convey("When I wish to publish things to google Pubsub", t, func() {
		expectedConfig := google.Config{
			ProjectID:       "testProject",
			CredentialsPath: "path/to/cred.yml",
			Concurrency:     150,
			Timeout:         666 * time.Second,
		}

		src, cli := initTestClient("something")
		defer src.Close()

		ps := google.Pubsub{
			Client: cli,
			Config: expectedConfig,
		}

		Convey("I should be able to initialize a Publish implementation", func() {
			publisher := ps.Publish()
			So(publisher, ShouldNotBeZeroValue)
		})
	})
}

func TestPubsub_Listen(t *testing.T) {

}

func TestPubsub_Receive(t *testing.T) {

}

func TestPubsub_Registry(t *testing.T) {

}

func ExampleNewPubsub_withoutOption() {
	conf := google.Config{
		ProjectID:       "aSuperCoolProject",
		CredentialsPath: "path/to/credentials.json",
		Timeout:         10 * time.Second,
		Concurrency:     10,
	}

	_, err := google.NewPubsub(context.Background(), conf)
	if err != nil {
		panic(err)
	}
}

func ExampleNewPubsub_withOption() {
	conf := google.Config{
		ProjectID:       "aSuperCoolProject",
		CredentialsPath: "path/to/credentials.json",
	}

	_, err := google.NewPubsub(
		context.Background(), conf,
		option.WithCredentialsFile("some/other/credentials.json"),
		option.WithEndpoint("dont.evil/know/where"),
	)
	if err != nil {
		panic(err)
	}

	// in this example, client uses credentials path from Config. Passing an option will not override
	// credentials values.
}
