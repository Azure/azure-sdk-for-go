package storage

import (
	"time"

	chk "gopkg.in/check.v1"
)

type StorageQueueSuite struct{}

var _ = chk.Suite(&StorageQueueSuite{})

func getQueueClient(c *chk.C) *QueueServiceClient {
	cli := getBasicClient(c).GetQueueService()
	return &cli
}

func (s *StorageQueueSuite) Test_pathForQueue(c *chk.C) {
	c.Assert(getQueueClient(c).
		GetQueueReference("q").
		buildPath(), chk.Equals, "/q")
}

func (s *StorageQueueSuite) Test_pathForQueueMessages(c *chk.C) {
	c.Assert(getQueueClient(c).
		GetQueueReference("q").
		buildPathMessages(), chk.Equals, "/q/messages")
}

func (s *StorageQueueSuite) TestCreateQueue_DeleteQueue(c *chk.C) {
	q := getQueueClient(c).GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	c.Assert(q.Delete(nil), chk.IsNil)
}

func (s *StorageQueueSuite) Test_GetMetadata_GetApproximateCount(c *chk.C) {
	cli := getQueueClient(c)
	q := cli.GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	err := q.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(q.AproxMessageCount, chk.Equals, uint64(0))

	m := q.GetMessageReference("lolrofl")
	for ix := 0; ix < 3; ix++ {
		err = m.Put(nil)
		c.Assert(err, chk.IsNil)
	}
	time.Sleep(1 * time.Second)

	err = q.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(q.AproxMessageCount, chk.Equals, uint64(3))
}

func (s *StorageQueueSuite) Test_SetAndGetMetadata(c *chk.C) {
	q := getQueueClient(c).GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	metadata := map[string]string{
		"Lol1":   "rofl1",
		"lolBaz": "rofl",
	}
	q.Metadata = metadata
	err := q.SetMetadata(nil)
	c.Assert(err, chk.IsNil)

	err = q.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(q.Metadata["lol1"], chk.Equals, metadata["Lol1"])
	c.Assert(q.Metadata["lolbaz"], chk.Equals, metadata["lolBaz"])
}

func (s *StorageQueueSuite) TestQueueExists(c *chk.C) {
	q := getQueueClient(c).GetQueueReference("nonexistent-queue")
	ok, err := q.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, false)

	q.Name = randString(20)
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	ok, err = q.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, true)
}

func (s *StorageQueueSuite) TestGetMessages(c *chk.C) {
	cli := getQueueClient(c)
	q := cli.GetQueueReference(randString(20))
	c.Assert(q.Create(nil), chk.IsNil)
	defer q.Delete(nil)

	m := q.GetMessageReference("message")
	n := 4
	for i := 0; i < n; i++ {
		c.Assert(m.Put(nil), chk.IsNil)
	}

	options := GetMessagesOptions{
		NumOfMessages: n,
	}
	list, err := q.GetMessages(&options)
	c.Assert(err, chk.IsNil)
	c.Assert(list, chk.HasLen, n)
}
