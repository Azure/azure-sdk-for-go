package storage

import chk "gopkg.in/check.v1"

type StorageMessageSuite struct{}

var _ = chk.Suite(&StorageMessageSuite{})

func (s *StorageMessageSuite) Test_pathForMessage(c *chk.C) {
	m := getQueueClient(c).GetQueueReference("q").GetMessageReference("m")
	m.ID = "ID"
	c.Assert(m.buildPath(), chk.Equals, "/q/messages/ID")
}

func (s *StorageMessageSuite) TestDeleteMessages(c *chk.C) {
	cli := getQueueClient(c)
	q := cli.GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	m := q.GetMessageReference("message")
	c.Assert(m.Put(nil), chk.IsNil)

	options := GetMessagesOptions{
		VisibilityTimeout: 1,
	}
	list, err := q.GetMessages(&options)
	c.Assert(err, chk.IsNil)
	c.Assert(list, chk.HasLen, 1)

	c.Assert(list[0].Delete(nil), chk.IsNil)
}

func (s *StorageMessageSuite) TestPutMessage_PeekMessage_UpdateMessage_DeleteMessage(c *chk.C) {
	cli := getQueueClient(c)
	q := cli.GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	m := q.GetMessageReference(randString(64 * 1024)) // exercise max length
	c.Assert(m.Put(nil), chk.IsNil)

	list, err := q.PeekMessages(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(list, chk.HasLen, 1)
	c.Assert(list[0].Text, chk.Equals, m.Text)

	getOptions := GetMessagesOptions{
		NumOfMessages:     1,
		VisibilityTimeout: 2,
	}
	list, gerr := q.GetMessages(&getOptions)
	c.Assert(gerr, chk.IsNil)

	m = &(list[0])
	m.Text = "Test Message"

	updateOptions := UpdateMessageOptions{
		VisibilityTimeout: 2,
	}
	c.Assert(m.Update(&updateOptions), chk.IsNil)

	list, err = q.PeekMessages(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(list, chk.HasLen, 0)
}
