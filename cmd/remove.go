package cmd

import (
	"errors"
	"math/rand"
	"time"

	"github.com/d-kuro/escale/pkg/aws"
	"github.com/d-kuro/escale/pkg/aws/ec2"
	"github.com/d-kuro/escale/pkg/elasticsearch"
	"github.com/d-kuro/escale/pkg/log"

	"github.com/spf13/cobra"
)

type RemoveOption struct {
	host             string
	port             int
	autoScalingGroup string
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
	fset.StringVarP(&o.autoScalingGroup, "auto-scaling-group", "g", "", "Auto Scaling Group name")
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

	target := dataNodes[rand.Intn(len(dataNodes))]
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
