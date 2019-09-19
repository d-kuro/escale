package autoscaling

import "github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"

type Client struct {
	api autoscalingiface.AutoScalingAPI
}

func NewClient(api autoscalingiface.AutoScalingAPI) *Client {
	return &Client{api: api}
}
