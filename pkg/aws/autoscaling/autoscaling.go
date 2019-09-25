package autoscaling

import (
	"fmt"

	"github.com/d-kuro/escale/pkg/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type Client struct {
	api autoscalingiface.AutoScalingAPI
}

func NewClient(provider client.ConfigProvider) *Client {
	api := autoscaling.New(provider)
	return &Client{api: api}
}

func (c *Client) IncreaseInstances(groupName string, desired int64) error {
	describeReq := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(groupName),
		},
	}
	resp, err := c.api.DescribeAutoScalingGroups(describeReq)
	if err != nil {
		return err
	}

	if len(resp.AutoScalingGroups) == 0 {
		return err
	}
	asg := resp.AutoScalingGroups[0]

	currentDesiredCapacity := aws.Int64Value(asg.DesiredCapacity)
	if currentDesiredCapacity >= desired {
		return fmt.Errorf(
			"current desired capacity is %d, please specify a larger number for desired capacity",
			currentDesiredCapacity)
	}

	log.Logger.Printf("AutoScalingGroup: %s, set desired capacity: %d -> %d",
		groupName, currentDesiredCapacity, desired)
	setCapReq := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(groupName),
		DesiredCapacity:      aws.Int64(desired),
	}
	if _, err = c.api.SetDesiredCapacity(setCapReq); err != nil {
		return err
	}

	return nil
}

func (c *Client) DetachInstance(groupName string, instanceID string) error {
	req := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName: aws.String(groupName),
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	}
	if _, err := c.api.DetachInstances(req); err != nil {
		return err
	}
	return nil
}
