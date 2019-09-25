package cmd

import (
	"errors"
	"math/rand"
	"time"

	"github.com/d-kuro/escale/pkg/aws"
	"github.com/d-kuro/escale/pkg/aws/autoscaling"
	"github.com/d-kuro/escale/pkg/aws/ec2"
	"github.com/d-kuro/escale/pkg/elasticsearch"
	"github.com/d-kuro/escale/pkg/log"

	"github.com/cenkalti/backoff"
	"github.com/github/vulcanizer"
	"github.com/spf13/cobra"
)

type RemoveOption struct {
	host             string
	port             int
	removeNodeName   string
	autoScalingGroup string
	maxRetry         uint64
	delay            time.Duration
	region           string
	profile          string
	configFile       string
}

func NewRemoveCommand(o *RemoveOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove Elasticsearch nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemoveCommand(o)
		},
	}
	fset := cmd.Flags()
	fset.StringVar(&o.host, "host", "", "Elasticsearch host")
	fset.IntVar(&o.port, "port", 9200, "Elasticsearch port")
	fset.StringVar(&o.removeNodeName, "remove-node-name", "", "The name of the node to remove")
	fset.StringVarP(&o.autoScalingGroup, "auto-scaling-group", "g", "", "Auto Scaling Group name")
	fset.Uint64Var(&o.maxRetry, "max-retry", 300, "Max retry count")
	fset.DurationVar(&o.delay, "delay", 5*time.Second, "Delay between retries")
	fset.StringVar(&o.region, "region", "", "AWS region")
	fset.StringVar(&o.profile, "profile", "", "AWS profile name")
	fset.StringVarP(&o.configFile, "config-file", "f", ".escale.yaml", "Configuration file to read in")
	return cmd
}

func runRemoveCommand(o *RemoveOption) error {
	rand.Seed(time.Now().UnixNano())

	o, err := o.readConfig()
	if err != nil {
		return err
	}
	if err := o.validate(); err != nil {
		return err
	}

	log.Logger.Printf("List Elasticsearch data nodes")
	esClient := elasticsearch.NewClient(o.host, o.port)
	nodes, err := esClient.ListNodes()
	if err != nil {
		return err
	}
	dataNodes := elasticsearch.ListDataNodes(nodes)
	for _, dataNode := range dataNodes {
		log.Logger.Printf("%+v", dataNode)
	}

	var target vulcanizer.Node
	if o.removeNodeName != "" {
		target, err = elasticsearch.GetNodeFromName(o.removeNodeName, dataNodes)
		if err != nil {
			return err
		}
	} else {
		target = dataNodes[rand.Intn(len(dataNodes))]
	}
	log.Logger.Printf("Remove target data node: %+v", target)

	sess, err := aws.NewSession(o.region, o.profile)
	if err != nil {
		return err
	}
	ec2Client := ec2.NewClient(sess)
	instanceID, err := ec2Client.RetrieveInstanceIDFromPrivateIP(target.Ip)
	if err != nil {
		return err
	}
	log.Logger.Printf("Remove target instanceID: %s", instanceID)

	if err := drainShardsFromNodes(target.Name, esClient, o); err != nil {
		return err
	}

	log.Logger.Printf("Detaching target instance from Auto Scaling Group: %s", o.autoScalingGroup)
	asgClient := autoscaling.NewClient(sess)
	if err := asgClient.DetachInstance(o.autoScalingGroup, instanceID); err != nil {
		return err
	}

	log.Logger.Printf("Terminate instance: %s", instanceID)
	if err := ec2Client.TerminateInstance(instanceID); err != nil {
		return err
	}

	log.Logger.Printf("Removes all shard allocation exclusion rules")
	if _, err := esClient.FillAll(); err != nil {
		return err
	}

	log.Logger.Printf("Finished")
	return nil
}

func (o RemoveOption) readConfig() (*RemoveOption, error) {
	if o.configFile == "" {
		return &o, nil
	}
	config, err := GetConfig(o.configFile)
	if err != nil {
		switch err.(type) {
		case *ErrFileNotExist:
			return &o, nil
		default:
			return nil, err
		}
	}

	if o.host == "" {
		o.host = config.Host
	}
	if o.port == 9200 {
		o.port = config.Port
	}
	if o.autoScalingGroup == "" {
		o.autoScalingGroup = config.AutoScalingGroup
	}
	return &o, nil
}

func (o RemoveOption) validate() error {
	if o.host == "" {
		return errors.New("host is required")
	}

	if o.autoScalingGroup == "" {
		return errors.New("auto-scaling-group is required")
	}
	return nil
}

func drainShardsFromNodes(nodeName string, client *elasticsearch.Client, o *RemoveOption) error {
	log.Logger.Printf("Drain shards from target node: %s", nodeName)
	if _, err := client.DrainServer(nodeName); err != nil {
		return err
	}

	log.Logger.Printf("Waiting for shards drain from target node...")
	backOff := backoff.WithMaxRetries(backoff.NewConstantBackOff(o.delay), o.maxRetry)
	if err := backoff.Retry(getShardsFromNodes([]string{nodeName}, client), backOff); err != nil {
		return err
	}
	return nil
}

func getShardsFromNodes(nodes []string, client *elasticsearch.Client) func() error {
	return func() error {
		nodes, err := client.GetShards(nodes)
		if err != nil {
			return err
		}
		if len(nodes) == 0 {
			return nil
		}
		return errors.New("timed out: shards do not drain from the given node")
	}
}

func listDataNodes() error {

}
