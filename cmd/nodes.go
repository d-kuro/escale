package cmd

import (
	"errors"

	"github.com/d-kuro/escale/pkg/elasticsearch"

	"github.com/d-kuro/escale/pkg/log"
	"github.com/spf13/cobra"
)

type NodesOption struct {
	host       string
	port       int
	configFile string
}

func NewNodesCommand(o *NodesOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "List Elasticsearch nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNodesCommand(o)
		},
	}
	fset := cmd.Flags()
	fset.StringVar(&o.host, "host", "", "Elasticsearch host")
	fset.IntVar(&o.port, "port", 9200, "Elasticsearch port")
	fset.StringVarP(&o.configFile, "config-file", "f", ".escale.yaml", "Configuration file to read in")
	return cmd
}

func runNodesCommand(o *NodesOption) error {
	o, err := o.readConfig()
	if err != nil {
		return err
	}
	if err := o.validate(); err != nil {
		return err
	}

	client := elasticsearch.NewClient(o.host, o.port)
	nodes, err := client.GetNodes()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		log.Logger.Printf("%+v", node)
	}
	return nil
}

func (o NodesOption) readConfig() (*NodesOption, error) {
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
	return &o, nil
}

func (o NodesOption) validate() error {
	if o.host == "" {
		return errors.New("host is required")
	}
	return nil
}
