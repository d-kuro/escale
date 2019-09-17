package cmd

import (
	"github.com/d-kuro/escale/pkg/log"
	"github.com/github/vulcanizer"
	"github.com/spf13/cobra"
)

type NodesOption struct {
	host   string
	port   int
	config string
}

func NewNodesCommand(o *NodesOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "List Elasticsearch nodes",
		Run: func(cmd *cobra.Command, args []string) {
			runNodesCommand(o)
		},
	}
	fset := cmd.Flags()
	fset.StringVar(&o.host, "host", "", "Elasticsearch host")
	fset.IntVar(&o.port, "port", 9200, "Elasticsearch port")
	fset.StringVarP(&o.host, "config", "f", ".escale.yaml", "Configuration file to read in")
	return cmd
}

func runNodesCommand(o *NodesOption) error {
	o, err := o.readConfig()
	if err != nil {
		return err
	}

	v := vulcanizer.NewClient(o.host, o.port)
	nodes, err := v.GetNodes()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		log.Logger.Printf("%+v", node)
	}
	return nil
}

func (o NodesOption) readConfig() (*NodesOption, error) {
	if o.config == "" {
		return &o, nil
	}
	config, err := GetConfig(o.config)
	if err != nil {
		return nil, err
	}
	if o.host == "" {
		o.host = config.Host
	}
	if o.port == 9200 {
		o.port = config.Port
	}
	return &o, nil
}
