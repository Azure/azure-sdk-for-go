package storage

import chk "gopkg.in/check.v1"

type LeaseBlobSuite struct{}

var _ = chk.Suite(&LeaseBlobSuite{})

func (s *LeaseBlobSuite) TestAcquireLeaseWithNoProposedLeaseID(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	_, err := b.AcquireLease(30, "")
	c.Assert(err, chk.IsNil)
}

func (s *LeaseBlobSuite) TestAcquireLeaseWithProposedLeaseID(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	leaseID, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)
	c.Assert(leaseID, chk.Equals, proposedLeaseID)
}

func (s *LeaseBlobSuite) TestAcquireLeaseWithBadProposedLeaseID(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	proposedLeaseID := "badbadbad"
	_, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.NotNil)
}

func (s *LeaseBlobSuite) TestRenewLeaseSuccessful(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	leaseID, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	err = b.RenewLease(leaseID)
	c.Assert(err, chk.IsNil)
}

func (s *LeaseBlobSuite) TestRenewLeaseAgainstNoCurrentLease(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	badLeaseID := "Golang rocks on Azure"
	err := b.RenewLease(badLeaseID)
	c.Assert(err, chk.NotNil)
}

func (s *LeaseBlobSuite) TestChangeLeaseSuccessful(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)
	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	leaseID, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	newProposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fbb"
	newLeaseID, err := b.ChangeLease(leaseID, newProposedLeaseID)
	c.Assert(err, chk.IsNil)
	c.Assert(newLeaseID, chk.Equals, newProposedLeaseID)
}

func (s *LeaseBlobSuite) TestChangeLeaseNotSuccessfulbadProposedLeaseID(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)
	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	leaseID, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	newProposedLeaseID := "1f812371-a41d-49e6-b123-f4b542e"
	_, err = b.ChangeLease(leaseID, newProposedLeaseID)
	c.Assert(err, chk.NotNil)
}

func (s *LeaseBlobSuite) TestReleaseLeaseSuccessful(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)
	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	leaseID, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	err = b.ReleaseLease(leaseID)
	c.Assert(err, chk.IsNil)
}

func (s *LeaseBlobSuite) TestReleaseLeaseNotSuccessfulBadLeaseID(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)
	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	_, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	err = b.ReleaseLease("badleaseid")
	c.Assert(err, chk.NotNil)
}

func (s *LeaseBlobSuite) TestBreakLeaseSuccessful(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	proposedLeaseID := "dfe6dde8-68d5-4910-9248-c97c61768fea"
	_, err := b.AcquireLease(30, proposedLeaseID)
	c.Assert(err, chk.IsNil)

	_, err = b.BreakLease()
	c.Assert(err, chk.IsNil)
}
