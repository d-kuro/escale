package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Client struct {
	api ec2iface.EC2API
}

func NewClient(provider client.ConfigProvider) *Client {
	api := ec2.New(provider)
	return &Client{api: api}
}

func (c *Client) RetrieveInstanceIDFromPrivateIP(privateIP string) (string, error) {
	resp, err := c.api.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("private-ip-address"),
				Values: []*string{
					aws.String(privateIP),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("instance with %q not found", privateIP)
	}

	return aws.StringValue(resp.Reservations[0].Instances[0].InstanceId), nil
}
