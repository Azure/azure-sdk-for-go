// Package cosmosdb provides clients for Microsoft Azure CosmosDb Services.
package cosmosdb

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	cosmosdb "github.com/Azure/azure-sdk-for-go/arm/cosmos-db"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"gopkg.in/mgo.v2"
)

const (
	cosmosDbConnectionPort = 10255
	managementResource     = "https://management.azure.com/"
)

// NewMongoDBClientWithConnectionString returns a MongoDb session to communicate with CosmosDB using a connection string.
func NewMongoDBClientWithConnectionString(connectionString string) (*mgo.Session, error) {

	// strip out the "ssl=true" option as MongoDb driver does not support by default SSL.
	connectionString = strings.Replace(connectionString, "ssl=true", "", -1)
	dialInfo, err := mgo.ParseURL(connectionString)

	if err != nil {
		return nil, err
	}

	return NewMongoDBClient(dialInfo)
}

// NewMongoDBClientWithCredentials returns a MongoDb session to communicate with CosmosDB using a username and a password.
func NewMongoDBClientWithCredentials(username, password, host string) (*mgo.Session, error) {

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%d", host, cosmosDbConnectionPort)},
		Username: username,
		Password: password,
	}

	return NewMongoDBClient(dialInfo)
}

// NewMongoDBClientWithSPToken returns a  session to communicate with CosmosDB using an auth token.
func NewMongoDBClientWithSPToken(spToken *adal.ServicePrincipalToken, subscriptionID, resourceGroup, account string) (*mgo.Session, error) {

	authorizer := autorest.NewBearerAuthorizer(spToken)

	cosmosDbClient := cosmosdb.NewDatabaseAccountsClient(subscriptionID)
	cosmosDbClient.Authorizer = authorizer

	result, err := cosmosDbClient.ListConnectionStrings(resourceGroup, account)

	if err != nil {
		return nil, err
	}

	connectionStrings := *result.ConnectionStrings

	fmt.Println(fmt.Sprintf("this isthe conn string: %s", *(connectionStrings[0].ConnectionString)))

	for _, connectionString := range connectionStrings {
		session, err := NewMongoDBClientWithConnectionString(*connectionString.ConnectionString)

		if session != nil && err == nil {
			return session, nil
		}
	}

	return nil, err
}

// NewMongoDBClientWithMSI returns a MongoDB session to communicate with CosmosDB using MSI.
func NewMongoDBClientWithMSI(subscriptionID, resourceGroup, account string) (*mgo.Session, error) {

	msiEndpoint, err := adal.GetMSIVMEndpoint()
	fmt.Println(msiEndpoint)

	spToken, err := adal.NewServicePrincipalTokenFromMSI(msiEndpoint, managementResource)

	if err != nil {
		return nil, err
	}

	return NewMongoDBClientWithSPToken(spToken, subscriptionID, resourceGroup, account)
}

// NewMongoDBClient returns a MongoDB session to communicate with CosmosDB.
func NewMongoDBClient(dialInfo *mgo.DialInfo) (*mgo.Session, error) {

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), &tls.Config{})
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return nil, err
	}

	return session, nil
}
