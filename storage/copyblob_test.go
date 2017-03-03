package storage

import (
	"io/ioutil"
	"net/http"
	"testing"

	chk "gopkg.in/check.v1"
)

type CopyBlobSuite struct{}

var _ = chk.Suite(&CopyBlobSuite{})

func (s *CopyBlobSuite) TestBlobCopy(c *chk.C) {
	if testing.Short() {
		c.Skip("skipping blob copy in short mode, no SLA on async operation")
	}

	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	srcBlob := cnt.GetBlobReference(randName(5))
	dstBlob := cnt.GetBlobReference(randName(5))
	body := []byte(randString(1024))

	c.Assert(srcBlob.putSingleBlockBlob(body), chk.IsNil)
	defer srcBlob.Delete(nil)

	c.Assert(dstBlob.Copy(srcBlob.GetURL()), chk.IsNil)
	defer dstBlob.Delete(nil)

	resp, err := dstBlob.Get()
	c.Assert(err, chk.IsNil)

	b, err := ioutil.ReadAll(resp)
	defer resp.Close()
	c.Assert(err, chk.IsNil)
	c.Assert(b, chk.DeepEquals, body)
}

func (s *CopyBlobSuite) TestStartBlobCopy(c *chk.C) {
	if testing.Short() {
		c.Skip("skipping blob copy in short mode, no SLA on async operation")
	}

	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	srcBlob := cnt.GetBlobReference(randName(5))
	dstBlob := cnt.GetBlobReference(randName(5))
	body := []byte(randString(1024))

	c.Assert(srcBlob.putSingleBlockBlob(body), chk.IsNil)
	defer srcBlob.Delete(nil)

	// given we dont know when it will start, can we even test destination creation?
	// will just test that an error wasn't thrown for now.
	copyID, err := dstBlob.StartCopy(srcBlob.GetURL())
	c.Assert(copyID, chk.NotNil)
	c.Assert(err, chk.IsNil)
}

// Tests abort of blobcopy. Given the blobcopy is usually over before we can actually trigger an abort
// it is agreed that we perform a copy then try and perform an abort. It should result in a HTTP status of 409.
// So basically we're testing negative scenario (as good as we can do for now)
func (s *CopyBlobSuite) TestAbortBlobCopy(c *chk.C) {
	if testing.Short() {
		c.Skip("skipping blob copy in short mode, no SLA on async operation")
	}

	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	srcBlob := cnt.GetBlobReference(randName(5))
	dstBlob := cnt.GetBlobReference(randName(5))
	body := []byte(randString(1024))

	c.Assert(srcBlob.putSingleBlockBlob(body), chk.IsNil)
	defer srcBlob.Delete(nil)

	// given we dont know when it will start, can we even test destination creation?
	// will just test that an error wasn't thrown for now.
	copyID, err := dstBlob.StartCopy(srcBlob.GetURL())
	c.Assert(copyID, chk.NotNil)
	c.Assert(err, chk.IsNil)

	err = dstBlob.WaitForCopy(copyID)
	c.Assert(err, chk.IsNil)

	// abort abort abort, but we *know* its already completed.
	err = dstBlob.AbortCopy(copyID, "", 0)

	// abort should fail (over already)
	c.Assert(err.(AzureStorageServiceError).StatusCode, chk.Equals, http.StatusConflict)
}
