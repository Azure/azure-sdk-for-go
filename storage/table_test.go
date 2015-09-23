package storage

import (
	"encoding/base64"
	chk "gopkg.in/check.v1"
)

type TableClient struct{}

func getTableClient(c *chk.C) TableServiceClient {
	return getBasicClient(c).GetTableService()
}

func (s *StorageBlobSuite) Test_SharedKeyLite(c *chk.C) {
	cli := getTableClient(c)

	// override the accountKey and accountName
	// but make sure to reset when returning
	oldAK := cli.client.accountKey
	oldAN := cli.client.accountName

	defer func() {
		cli.client.accountKey = oldAK
		cli.client.accountName = oldAN
	}()

	// don't worry, I've already changed mine :)
	key, err := base64.StdEncoding.DecodeString("zHDHGs7C+Di9pZSDMuarxJJz3xRBzAHBYaobxpLEc7kwTptR/hPEa9j93hIfb2Tbe9IA50MViGmjQ6nUF/OVvA==")
	if err != nil {
		c.Fail()
	}

	cli.client.accountKey = key
	cli.client.accountName = "mindgotest"

	headers := map[string]string{
		"Accept-Charset": "UTF-8",
		"Content-Type":   "application/json",
		"x-ms-date":      "Wed, 23 Sep 2015 16:40:05 GMT",
		"Content-Length": "0",
		"x-ms-version":   "2015-02-21",
		"Accept":         "application/json;odata=nometadata",
	}
	url := "https://mindgotest.table.core.windows.net/tquery()"

	ret, err := cli.client.createSharedKeyLite(url, headers)
	if err != nil {
		c.Fail()
	}

	c.Assert(ret, chk.Equals, "SharedKeyLite mindgotest:+32DTgsPUgXPo/O7RYaTs0DllA6FTXMj3uK4Qst8y/E=")
}

func (s *StorageBlobSuite) Test_CreateTable(c *chk.C) {
	cli := getTableClient(c)

	err := cli.CreateTable("longtimeagoinagalaxyfarfaraway")

	c.Assert(err, chk.IsNil)
}
