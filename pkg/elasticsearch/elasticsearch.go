package elasticsearch

import "github.com/github/vulcanizer"

type Client struct {
	client *vulcanizer.Client
}

func NewClient(host string, port int) *Client {
	v := vulcanizer.NewClient(host, port)
	return &Client{client: v}
}

func (c *Client) GetNodes() ([]vulcanizer.Node, error) {
	return c.client.GetNodes()
}
