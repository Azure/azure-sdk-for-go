package applicationGateway

import (
	origin "github.com/Azure/azure-sdk-for-go/arm/network/2017-03-01/applicationGateway"
)

type ApplicationGatewaysClient struct {
	origin.ApplicationGatewaysClient
}

func (client ApplicationGatewaysClient) BackendHealth(resourceGroupName, applicationGatewayName, expand string, cancel <-chan struct{}) (<-chan BackendHealthType, <-chan error) {
	results, errs := make(chan BackendHealthType), make(chan error)

	go func() {
		defer close(results)
		defer close(errs)

		intermediateResults, intermediateErrs := client.ApplicationGatewaysClient.BackendHealth(resourceGroupName, applicationGatewayName, expand, cancel)

		results <- BackendHealthType(<-intermediateResults)
		errs <- <-intermediateErrs
	}()

	return results, errs
}
