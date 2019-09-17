package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type AddOption struct {
	host             string
	port             int
	autoScalingGroup string
	number           int
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
	fset.IntVarP(&o.number, "number", "n", 1, "Number to add instances")
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
	return nil
}
