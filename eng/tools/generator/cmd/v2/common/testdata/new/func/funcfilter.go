package funcfilter_test

type Client struct{}

func (client *Client) BeginCreateOrUpdate(resourceGroupName string, options *ClientBeginCreateOrUpdateOptions) (ClientBeginCreateOrUpdateResponse, error) {

	return ClientBeginCreateOrUpdateResponse{}, nil
}

type ClientBeginCreateOrUpdateOptions struct{}

type ClientBeginCreateOrUpdateResponse struct{}
