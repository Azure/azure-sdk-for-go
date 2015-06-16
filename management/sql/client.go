package sql

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/management"
)

// Definitions of numerous constants representing API endpoints.
const (
	azureCreateDatabaseServerURL = "services/sqlservers/servers"
	azureListDatabaseServersURL  = "/services/sqlservers/servers"
	azureDeleteDatabaseServerURL = "/services/sqlservers/servers/%s"

	azureCreateDatabaseURL = "services/sqlservers/servers/%s/databases"
	azureGetDatabaseURL    = "/services/sqlservers/servers/%s/databases/%s"
	azureListDatabasesURL  = "/services/sqlservers/servers/%s/databases?contentview=generic"
	azureUpdateDatabaseURL = "/services/sqlservers/servers/%s/databases/%s"
	azureDeleteDatabaseURL = "/services/sqlservers/servers/%s/databases/%s"

	errParamNotSpecified = "Parameter %s was not specified."

	DatabaseStateCreating = "Creating"
)

// SqlDatabaseClient defines various database CRUD operations.
// It contains a management.Client for making the actual http calls.
type SqlDatabaseClient struct {
	mgmtClient management.Client
}

// NewClient returns a new SqlDatabaseClient struct with the provided
// management.Client as the underlying client.
func NewClient(mgmtClient management.Client) SqlDatabaseClient {
	return SqlDatabaseClient{mgmtClient}
}

// CreateServer creates a new Azure SQL Database server and return its name.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505699.aspx
func (c SqlDatabaseClient) CreateServer(params DatabaseServerCreateParams) (string, error) {
	req, err := xml.Marshal(params)
	if err != nil {
		return "", err
	}

	resp, err := c.mgmtClient.SendAzurePostRequestWithReturnedResponse(azureCreateDatabaseServerURL, req)
	if err != nil {
		return "", err
	}

	var name string
	err = xml.Unmarshal(resp, &name)

	return name, err
}

// ListServers retrieves the Azure SQL Database servers for this subscription.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505702.aspx
func (c SqlDatabaseClient) ListServers() (ListServersResponse, error) {
	var resp ListServersResponse

	data, err := c.mgmtClient.SendAzureGetRequest(azureListDatabaseServersURL)
	if err != nil {
		return resp, err
	}

	err = xml.Unmarshal(data, &resp)
	return resp, err
}

// DeleteServer deletes an Azure SQL Database server (including all its databases).
//
// https://msdn.microsoft.com/en-us/library/azure/dn505695.aspx
func (c SqlDatabaseClient) DeleteServer(name string) error {
	if name == "" {
		return fmt.Errorf(errParamNotSpecified, "name")
	}

	url := fmt.Sprintf(azureDeleteDatabaseServerURL, name)
	_, err := c.mgmtClient.SendAzureDeleteRequest(url)
	return err
}

// CreateDatabase creates a new Microsoft Azure SQL Database on the given database server.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505701.aspx
func (c SqlDatabaseClient) CreateDatabase(server string, params DatabaseCreateParams) error {
	if server == "" {
		return fmt.Errorf(errParamNotSpecified, "server")
	}

	req, err := xml.Marshal(params)
	if err != nil {
		return err
	}

	target := fmt.Sprintf(azureCreateDatabaseURL, server)
	_, err = c.mgmtClient.SendAzurePostRequest(target, req)
	return err
}

// WaitForDatabaseCreation is a helper method which waits
// for the creation of the database on the given server.
func (c SqlDatabaseClient) WaitForDatabaseCreation(
	server, database string,
	cancel chan struct{}) error {
	for {
		stat, err := c.GetDatabase(server, database)
		if err != nil {
			return err
		}
		if stat.State != DatabaseStateCreating {
			return nil
		}

		select {
		case <-time.After(management.DefaultOperationPollInterval):
		case <-cancel:
			return management.ErrOperationCancelled
		}
	}
}

// GetDatabase gets the details for an Azure SQL Database.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505708.aspx
func (c SqlDatabaseClient) GetDatabase(server, database string) (ServiceResource, error) {
	var db ServiceResource

	if database == "" {
		return db, fmt.Errorf(errParamNotSpecified, "database")
	}
	if server == "" {
		return db, fmt.Errorf(errParamNotSpecified, "server")
	}

	url := fmt.Sprintf(azureGetDatabaseURL, server, database)
	resp, err := c.mgmtClient.SendAzureGetRequest(url)
	if err != nil {
		return db, err
	}

	err = xml.Unmarshal(resp, &db)
	return db, err
}

// ListDatabases returns a list of Azure SQL Databases on the given server.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505711.aspx
func (c SqlDatabaseClient) ListDatabases(server string) (ListDatabasesResponse, error) {
	var databases ListDatabasesResponse
	if server == "" {
		return databases, fmt.Errorf(errParamNotSpecified, "server name")
	}

	url := fmt.Sprintf(azureListDatabasesURL, server)
	resp, err := c.mgmtClient.SendAzureGetRequest(url)
	if err != nil {
		return databases, err
	}

	err = xml.Unmarshal(resp, &databases)
	return databases, err
}

// UpdateDatabase updates the details of the given Database off the given server.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505718.aspx
func (c SqlDatabaseClient) UpdateDatabase(
	server, database string,
	params ServiceResourceUpdateParams) (management.OperationID, error) {
	if database == "" {
		return "", fmt.Errorf(errParamNotSpecified, "database")
	}
	if server == "" {
		return "", fmt.Errorf(errParamNotSpecified, "server")
	}

	url := fmt.Sprintf(azureUpdateDatabaseURL, server, database)
	req, err := xml.Marshal(params)
	if err != nil {
		return "", err
	}

	return c.mgmtClient.SendAzurePutRequest(url, "text/xml", req)
}

// DeleteDatabase deletes the Azure SQL Database off the given server.
//
// https://msdn.microsoft.com/en-us/library/azure/dn505705.aspx
func (c SqlDatabaseClient) DeleteDatabase(server, database string) error {
	if database == "" {
		return fmt.Errorf(errParamNotSpecified, "database")
	}
	if server == "" {
		return fmt.Errorf(errParamNotSpecified, "server")
	}

	url := fmt.Sprintf(azureDeleteDatabaseURL, server, database)

	_, err := c.mgmtClient.SendAzureDeleteRequest(url)

	return err
}
