package funcfilter_test

type Client struct{}

func (client *Client) Update(resourceGroupName string, options *ClientUpdateOptions) (ClientUpdateResponse, error) {

	return ClientUpdateResponse{}, nil
}

type ClientUpdateOptions struct{}

type ClientUpdateResponse struct{}
