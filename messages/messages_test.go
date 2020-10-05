package messages_test

// func TestBaseMessage(t *testing.T) {
// 	Convey("We should be able to init V1 messages", t, func() {
// 		So(messages.NewOtherV1("test", uuid.New()), ShouldNotBeZeroValue)
// 		So(messages.NewEmptyOtherV1(), ShouldNotBeNil)
// 	})
// }
//
// func TestMock(t *testing.T) {
// 	Convey("We should be able to", t, func() {
// 		m := messages.NewMock()
// 		te := &test{}
//
// 		Convey("mock ToPubsubMessage", func() {
// 			args := map[string]string{"test": "test"}
// 			data := []byte("testazerty")
// 			m.On("ToPubsubMessage").Return(args, data, errors.New("test"))
// 			actArgs, actData, err := m.ToPubsubMessage()
// 			So(actArgs, ShouldResemble, args)
// 			So(actData, ShouldResemble, data)
// 			So(err, ShouldResemble, errors.New("test"))
// 		})
//
// 		Convey("mock FromPubsubMessage", func() {
// 			args := map[string]string{"test": "test"}
// 			data := []byte("testazerty")
// 			m.On("FromPubsubMessage", args, data, mock.Anything, mock.Anything).Return(errors.New("test"))
// 			So(m.FromPubsubMessage(args, data, te.Ack, te.Nack), ShouldResemble, errors.New("test"))
// 		})
//
// 		Convey("mock Ack", func() {
// 			args := map[string]string{"test": "test"}
// 			data := []byte("testazerty")
// 			m.On("FromPubsubMessage", args, data, mock.Anything, mock.Anything).Return(nil)
// 			_ = m.FromPubsubMessage(args, data, te.Ack, te.Nack)
// 			m.On("Ack")
// 			m.Ack()
// 			So(te.Item, ShouldEqual, "ack")
// 		})
//
// 		Convey("mock Nack", func() {
// 			args := map[string]string{"test": "test"}
// 			data := []byte("testazerty")
// 			m.On("FromPubsubMessage", args, data, mock.Anything, mock.Anything).Return(nil)
// 			_ = m.FromPubsubMessage(args, data, te.Ack, te.Nack)
// 			m.On("Nack")
// 			m.Nack()
// 			So(te.Item, ShouldEqual, "nack")
// 		})
//
// 		Convey("mock Version", func() {
// 			m.On("Version").Return("test")
// 			So(m.Version(), ShouldEqual, "test")
// 		})
//
// 		Convey("mock Type", func() {
// 			m.On("Type").Return("test")
// 			So(m.Type(), ShouldEqual, "test")
// 		})
// 		coreTest.ResetMock(m.GetMock())
// 	})
// }

type test struct {
	Item string
}

func (t *test) Ack() {
	t.Item = "ack"
}

func (t *test) Nack() {
	t.Item = "nack"
}
