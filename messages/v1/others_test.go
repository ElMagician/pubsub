package v1_test

// func TestNewOther(t *testing.T) {
// 	Convey("We should be able to create an Other typed message", t, func() {
// 		expectedOther := v1.Other{
// 			Payload: map[string]interface{}{"test": "test"},
// 		}
// 		expectedOther.Source = "src"
// 		expectedOther.EmittedAt = time.Now()
//
// 		event := v1.NewOther("src", map[string]interface{}{"test": "test"})
//
// 		So(*event, ShouldResemble, expectedOther)
// 	})
// }
//
// func TestEmptyNewOther(t *testing.T) {
// 	Convey("We should be able to create an empty Other typed message", t, func() {
// 		event := v1.NewEmptyOther()
// 		So(*event, ShouldResemble, v1.Other{})
// 	})
// }
//
// func TestOther(t *testing.T) {
// 	Convey("Given an Other type message we should", t, func() {
// 		expectedOther := v1.NewOther("src", map[string]interface{}{"test": "test"})
// 		expectedDataString := "{\"payload\":{\"test\":\"test\"}}"
//
// 		Convey("be able to convert it to pubsub message data", func() {
// 			metaData, data, err := expectedOther.ToPubsubMessage()
// 			So(
// 				metaData,
// 				ShouldResemble,
// 				map[string]string{
// 					"emittedAt": strconv.FormatInt(expectedOther.EmittedAt, 10),
// 					"version":   v1.VersionKey,
// 					"type":      v1.OtherType,
// 					"source":    "src",
// 				},
// 			)
// 			So(string(data), ShouldEqual, expectedDataString)
// 			So(err, ShouldBeNil)
// 		})
//
// 		Convey("errored if not marshable", func() {
// 			expectedOther.Payload = make(chan string, 5)
// 			_, _, err := expectedOther.ToPubsubMessage()
// 			So(err, ShouldNotBeNil)
// 		})
//
// 		Convey("be able to get version", func() {
// 			So(expectedOther.Version(), ShouldEqual, v1.VersionKey)
// 		})
//
// 		Convey("be able to get type", func() {
// 			So(expectedOther.Type(), ShouldEqual, v1.OtherType)
// 		})
//
// 		Convey("be able to retrieve ", func() {
// 			testS := testStruct{Test: "rng"}
//
// 			metaData := map[string]string{
// 				"emittedAt": "666",
// 				"version":   v1.VersionKey,
// 				"type":      v1.OtherType,
// 				"source":    "src",
// 			}
// 			object := []byte(expectedDataString)
//
// 			expectedOther = &v1.Other{}
// 			err := expectedOther.FromPubsubMessage(metaData, object, testS.Ack, testS.Nack)
// 			So(err, ShouldBeNil)
//
// 			So(expectedOther.Version(), ShouldEqual, v1.VersionKey)
// 			So(expectedOther.Type(), ShouldEqual, v1.OtherType)
//
// 			// Data
// 			So(expectedOther.Payload, ShouldResemble, map[string]interface{}{"test": "test"})
//
// 			// MetaData
// 			So(expectedOther.Source, ShouldEqual, "src")
// 			So(expectedOther.EmittedAt, ShouldEqual, int64(666))
//
// 			Convey("and ack message", func() {
// 				expectedOther.Ack()
// 				So(testS.Test, ShouldEqual, "Ack")
// 			})
// 			Convey("and nack message", func() {
// 				expectedOther.Nack()
// 				So(testS.Test, ShouldEqual, "Nack")
// 			})
// 		})
//
// 		Convey("errored if data is not a correctly formatted AuditPS JSON", func() {
// 			metaData := map[string]string{
// 				"version": v1.VersionKey,
// 				"type":    v1.AuditType,
// 				"source":  "src",
// 			}
// 			object := []byte("{\"content\":")
// 			expectedOther = &v1.Other{}
// 			err := expectedOther.FromPubsubMessage(metaData, object, func() {}, func() {})
// 			So(err, ShouldNotBeNil)
// 		})
//
// 	})
// }

type testStruct struct {
	Test string
}

func (t *testStruct) Ack() {
	t.Test = "Ack"
}

func (t *testStruct) Nack() {
	t.Test = "Nack"
}
