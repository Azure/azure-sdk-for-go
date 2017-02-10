package storage

import (
	"crypto/rand"
	"net/url"
	"time"

	chk "gopkg.in/check.v1"
)

type StorageTableSuite struct{}

var _ = chk.Suite(&StorageTableSuite{})

func getTableClient(c *chk.C) TableServiceClient {
	return getBasicClient(c).GetTableService()
}

func deleteAllTables(c *chk.C) {
	cli := getBasicClient(c).GetTableService()
	result, err := cli.QueryTables(url.Values{})
	c.Assert(err, chk.IsNil)

	for _, t := range result.Tables {
		err := t.Delete()
		c.Assert(err, chk.IsNil)
	}
}

func (s *StorageTableSuite) Test_CreateAndDeleteTable(c *chk.C) {
	cli := getBasicClient(c).GetTableService()
	table := cli.GetTableReference(randTable())

	err := table.Create(EmptyPayload)
	c.Assert(err, chk.IsNil)

	// update table metadata
	table2 := cli.GetTableReference(randTable())
	err = table2.Create(FullMetadata)
	defer table2.Delete()
	c.Assert(err, chk.IsNil)
	// Check not empty values
	c.Assert(table2.OdataEditLink, chk.Not(chk.Equals), "")
	c.Assert(table2.OdataID, chk.Not(chk.Equals), "")
	c.Assert(table2.OdataMetadata, chk.Not(chk.Equals), "")
	c.Assert(table2.OdataType, chk.Not(chk.Equals), "")

	err = table.Delete()
	c.Assert(err, chk.IsNil)
}

func (s *StorageTableSuite) TestQueryTablesNextResults(c *chk.C) {
	deleteAllTables(c)
	cli := getBasicClient(c).GetTableService()

	for i := 0; i < 3; i++ {
		table := cli.GetTableReference(randTable())
		err := table.Create(EmptyPayload)
		c.Assert(err, chk.IsNil)
		defer table.Delete()
	}

	result, err := cli.QueryTables(url.Values{"top": {"2"}})
	c.Assert(err, chk.IsNil)
	c.Assert(result.Tables, chk.HasLen, 2)
	c.Assert(result.NextLink, chk.NotNil)

	result, err = cli.QueryTablesNextResults(result)
	c.Assert(err, chk.IsNil)
	c.Assert(result.Tables, chk.HasLen, 1)
	c.Assert(result.NextLink, chk.IsNil)

	result, err = cli.QueryTablesNextResults(result)
	c.Assert(result, chk.IsNil)
	c.Assert(err, chk.NotNil)
}

func randTable() string {
	const alphanum = "abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func appendTablePermission(policies []TableAccessPolicy, ID string,
	canRead bool, canAppend bool, canUpdate bool, canDelete bool,
	startTime time.Time, expiryTime time.Time) []TableAccessPolicy {

	tap := TableAccessPolicy{
		ID:         ID,
		StartTime:  startTime,
		ExpiryTime: expiryTime,
		CanRead:    canRead,
		CanAppend:  canAppend,
		CanUpdate:  canUpdate,
		CanDelete:  canDelete,
	}
	policies = append(policies, tap)
	return policies
}

func (s *StorageTableSuite) TestSetPermissionsSuccessfully(c *chk.C) {
	cli := getTableClient(c)
	table := cli.GetTableReference(randTable())
	c.Assert(table.Create(EmptyPayload), chk.IsNil)
	defer table.Delete()

	policies := []TableAccessPolicy{}
	policies = appendTablePermission(policies, "GolangRocksOnAzure", true, true, true, true, now, now.Add(10*time.Hour))
	table.AccessPolicies = policies

	err := table.SetPermissions(0)
	c.Assert(err, chk.IsNil)
}

func (s *StorageTableSuite) TestSetPermissionsUnsuccessfully(c *chk.C) {
	cli := getTableClient(c)
	table := cli.GetTableReference("nonexistingtable")

	policies := []TableAccessPolicy{}
	policies = appendTablePermission(policies, "GolangRocksOnAzure", true, true, true, true, now, now.Add(10*time.Hour))
	table.AccessPolicies = policies

	err := table.SetPermissions(0)
	c.Assert(err, chk.NotNil)
}

func (s *StorageTableSuite) TestSetThenGetPermissionsSuccessfully(c *chk.C) {
	cli := getTableClient(c)
	table := cli.GetTableReference(randTable())
	c.Assert(table.Create(EmptyPayload), chk.IsNil)
	defer table.Delete()

	policies := []TableAccessPolicy{}
	policies = appendTablePermission(policies, "GolangRocksOnAzure", true, true, true, true, now, now.Add(10*time.Hour))
	policies = appendTablePermission(policies, "AutoRestIsSuperCool", true, true, false, true, now.Add(20*time.Hour), now.Add(30*time.Hour))
	table.AccessPolicies = policies

	err := table.SetPermissions(0)
	c.Assert(err, chk.IsNil)

	err = table.GetPermissions(0)
	c.Assert(err, chk.IsNil)

	// now check policy set.
	c.Assert(table.AccessPolicies, chk.HasLen, 2)

	for i := range table.AccessPolicies {
		c.Assert(table.AccessPolicies[i].ID, chk.Equals, policies[i].ID)

		// test timestamps down the second
		// rounding start/expiry time original perms since the returned perms would have been rounded.
		// so need rounded vs rounded.
		c.Assert(table.AccessPolicies[i].StartTime.UTC().Round(time.Second).Format(time.RFC1123),
			chk.Equals, policies[i].StartTime.UTC().Round(time.Second).Format(time.RFC1123))
		c.Assert(table.AccessPolicies[i].ExpiryTime.UTC().Round(time.Second).Format(time.RFC1123),
			chk.Equals, policies[i].ExpiryTime.UTC().Round(time.Second).Format(time.RFC1123))

		c.Assert(table.AccessPolicies[i].CanRead, chk.Equals, policies[i].CanRead)
		c.Assert(table.AccessPolicies[i].CanAppend, chk.Equals, policies[i].CanAppend)
		c.Assert(table.AccessPolicies[i].CanUpdate, chk.Equals, policies[i].CanUpdate)
		c.Assert(table.AccessPolicies[i].CanDelete, chk.Equals, policies[i].CanDelete)
	}
}
