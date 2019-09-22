package elasticsearch

import (
	"time"

	"github.com/github/vulcanizer"
)

type Client struct {
	client *vulcanizer.Client
}

func NewClient(host string, port int) *Client {
	v := vulcanizer.NewClient(host, port)
	v.Timeout = 5 * time.Second
	return &Client{client: v}
}

func (c *Client) ListNodes() ([]vulcanizer.Node, error) {
	return c.client.GetNodes()
}

func (c *Client) GetShardOverlap(nodes []string) (map[string]vulcanizer.ShardOverlap, error) {
	return c.client.GetShardOverlap(nodes)
}

func ListDataNodes(nodes []vulcanizer.Node) []vulcanizer.Node {
	dataNodes := make([]vulcanizer.Node, 0, len(nodes)-1)
	for _, node := range nodes {
		if node.Master != "*" {
			dataNodes = append(dataNodes, node)
		}
	}
	return dataNodes
}
