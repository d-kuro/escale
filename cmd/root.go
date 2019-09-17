package cmd

import (
	"github.com/d-kuro/escale/pkg/log"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

func Execute() int {
	log.NewStdLogger()

	cmd := NewRootCommand()
	addCommands(cmd)

	if err := cmd.Execute(); err != nil {
		log.Logger.Errorf("error: %v", err)
		return exitCodeErr
	}
	return exitCodeOK
}

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "escale",
		Short:         "Elasticsearch node controller with AWS Auto Scaling Group",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewVersionCommand(),
		NewNodesCommand(&NodesOption{}),
	)
}
