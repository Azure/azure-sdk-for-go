package lrofilter_test

type Client struct{}

func (client *Client) CreateOrUpdate(resourceGroupName string, options *ClientCreateOrUpdateOptions) (ClientCreateOrUpdateResponse, error) {

	return ClientCreateOrUpdateResponse{}, nil
}

type ClientCreateOrUpdateOptions struct{}

type ClientCreateOrUpdateResponse struct{}
