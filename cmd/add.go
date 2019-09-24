package cmd

import (
	"errors"
	"time"

	"github.com/d-kuro/escale/pkg/aws"
	"github.com/d-kuro/escale/pkg/aws/autoscaling"
	"github.com/d-kuro/escale/pkg/elasticsearch"
	"github.com/d-kuro/escale/pkg/log"

	"github.com/cenkalti/backoff"
	"github.com/spf13/cobra"
)

const maxRetry = 100

type AddOption struct {
	host             string
	port             int
	autoScalingGroup string
	desired          int
	region           string
	profile          string
	configFile       string
}

func NewAddCommand(o *AddOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add Elasticsearch nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddCommand(o)
		},
	}
	fset := cmd.Flags()
	fset.StringVar(&o.host, "host", "", "Elasticsearch host")
	fset.IntVar(&o.port, "port", 9200, "Elasticsearch port")
	fset.StringVarP(&o.autoScalingGroup, "auto-scaling-group", "g", "", "Auto Scaling Group name")
	fset.IntVar(&o.desired, "desired", 0, "Desired capacity")
	fset.StringVar(&o.region, "region", "", "AWS region")
	fset.StringVar(&o.profile, "profile", "", "AWS profile name")
	fset.StringVarP(&o.configFile, "config-file", "f", ".escale.yaml", "Configuration file to read in")
	return cmd
}

func runAddCommand(o *AddOption) error {
	o, err := o.readConfig()
	if err != nil {
		return err
	}
	if err := o.validate(); err != nil {
		return err
	}

	esClient := elasticsearch.NewClient(o.host, o.port)

	log.Logger.Printf("Disable allocation for the cluster")
	if _, err := esClient.SetAllocation(false); err != nil {
		return err
	}

	sess, err := aws.NewSession(o.region, o.profile)
	if err != nil {
		return err
	}
	client := autoscaling.NewClient(sess)
	if err := client.IncreaseInstances(o.autoScalingGroup, int64(o.desired)); err != nil {
		return err
	}

	log.Logger.Printf("Waiting for nodes join to Elasticsearch cluster...")
	backOff := backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), maxRetry)
	if err := backoff.Retry(getElasticsearchNodes(esClient, o), backOff); err != nil {
		return err
	}

	log.Logger.Printf("Enable allocation for the cluster")
	if _, err := esClient.SetAllocation(true); err != nil {
		return err
	}
	log.Logger.Printf("Finished")
	return nil
}

func (o AddOption) readConfig() (*AddOption, error) {
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

func (o AddOption) validate() error {
	if o.host == "" {
		return errors.New("host is required")
	}
	if o.autoScalingGroup == "" {
		return errors.New("auto-scaling-group is required")
	}
	if o.desired <= 0 {
		return errors.New("desired is required (desired > 0)")
	}
	return nil
}

func getElasticsearchNodes(client *elasticsearch.Client, o *AddOption) func() error {
	return func() error {
		nodes, err := client.ListNodes()
		if err != nil {
			return err
		}
		if len(nodes) == o.desired {
			return nil
		}
		return errors.New("timed out: added nodes do not join to Elasticsearch cluster")
	}
}
