package azblob

import (
	"sort"

	chk "gopkg.in/check.v1"
)

// tests general functionality
func (s *aztestsSuite) TestBlobListWrapper(c *chk.C) {
	bsu, err := getGenericBSU("")

	c.Assert(err, chk.IsNil)

	container, _ := getContainerClient(c, bsu)

	_, err = container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	defer container.Delete(ctx, nil)

	files := []string{"a", "b", "c"}

	createNewBlobs(c, container, files)

	blobs, errs := container.ListBlobsFlatSegment(ctx, 3, 0, nil)

	found := make([]string, 0)

	for blobItem := range blobs {
		found = append(found, *blobItem.Name)
	}
	c.Assert(<- errs, chk.IsNil)

	sort.Strings(files)
	sort.Strings(found)

	c.Assert(found, chk.DeepEquals, files)
}

// tests that the buffer filling isn't a problem
func (s *aztestsSuite) TestBlobListWrapperFullBuffer(c *chk.C) {
	bsu, err := getGenericBSU("")

	c.Assert(err, chk.IsNil)

	container, _ := getContainerClient(c, bsu)

	_, err = container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	defer container.Delete(ctx, nil)

	files := []string{"a", "b", "c"}

	createNewBlobs(c, container, files)

	blobs, errs := container.ListBlobsFlatSegment(ctx, 1, 0, nil)

	found := make([]string, 0)

	for blobItem := range blobs {
		found = append(found, *blobItem.Name)
	}
	c.Assert(<- errs, chk.IsNil)

	sort.Strings(files)
	sort.Strings(found)

	c.Assert(files, chk.DeepEquals, found)
}

func (s *aztestsSuite) TestBlobListWrapperListingError(c *chk.C) {
	bsu, err := getGenericBSU("")

	c.Assert(err, chk.IsNil)

	container, _ := getContainerClient(c, bsu)

	blobs, errs := container.ListBlobsFlatSegment(ctx, 1, 0, nil)

	for _ = range blobs {
		// there should be NO blob listings coming back. Just an error.
		c.FailNow()
	}

	err = <- errs
	c.Log(err.Error())
	c.Assert(err, chk.NotNil)
}
