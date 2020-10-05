package v1_test

// func TestFromPubsubMessage(t *testing.T) {
// 	Convey("Given V1 base info we should", t, func() {
// 		Convey("have an error when trying to recover meta data without emission date from pubsub", func() {
// 			item := &v1.V1{}
// 			err := item.FromPubsubMessage(map[string]string{"key": "data"}, func() {}, func() {})
// 			So(err, ShouldNotBeNil)
// 			So(errors.Is(err, v1.ErrInvalidEmittedAt), ShouldBeTrue)
// 		})
//
// 		Convey("have an error when trying to recover meta data with incorrect emission date from pubsub", func() {
// 			item := &v1.V1{}
// 			err := item.FromPubsubMessage(map[string]string{"emittedAt": "data"}, func() {}, func() {})
// 			So(err, ShouldNotBeNil)
// 			So(errors.Is(err, v1.ErrInvalidEmittedAt), ShouldBeTrue)
// 		})
// 	})
// }
