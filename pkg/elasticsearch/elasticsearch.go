package elasticsearch

import (
	"fmt"
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

func (c *Client) GetShards(nodes []string) ([]vulcanizer.Shard, error) {
	return c.client.GetShards(nodes)
}

func (c *Client) DrainServer(nodeName string) (vulcanizer.ExcludeSettings, error) {
	return c.client.DrainServer(nodeName)
}

func (c *Client) FillAll() (vulcanizer.ExcludeSettings, error) {
	return c.client.FillAll()
}

func (c *Client) SetAllocation(allocation bool) (string, error) {
	if allocation {
		return c.client.SetAllocation("enable")
	}
	return c.client.SetAllocation("disable")
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

func GetNodeFromName(name string, nodes []vulcanizer.Node) (vulcanizer.Node, error) {
	for _, node := range nodes {
		if node.Name == name {
			return node, nil
		}
	}
	return vulcanizer.Node{}, fmt.Errorf("node not found name: %s", name)
}
