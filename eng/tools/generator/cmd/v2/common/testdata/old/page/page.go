package page

type Client struct{}

func (client *Client) GetLog(resourceGroupName string, options *ClientGetLogOptions) (ClientGetLogResponse, error) {

	return ClientGetLogResponse{}, nil
}

type ClientGetLogOptions struct{}

type ClientGetLogResponse struct{}

func (client *Client) NewListPager(resourceGroupName string, options *ClientNewListPagerOptions) (ClientNewListPagerResponse, error) {

	return ClientNewListPagerResponse{}, nil
}

type ClientNewListPagerOptions struct{}

type ClientNewListPagerResponse struct{}
