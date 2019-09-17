package cmd

import (
	"io"

	"github.com/d-kuro/escale/pkg/log"

	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

type Option struct {
	outStream io.Writer
	errStream io.Writer
}

func Execute(outStream, errStream io.Writer) int {
	log.NewStdLogger()

	o := NewOption(outStream, errStream)
	cmd := NewRootCommand(o)
	addCommands(cmd, o)

	if err := cmd.Execute(); err != nil {
		log.Logger.Errorf("error: %v", err)
		return exitCodeErr
	}
	return exitCodeOK
}

func NewRootCommand(o *Option) *cobra.Command {
	return &cobra.Command{
		Use:           "escale",
		Short:         "",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

func NewOption(outStream, errStream io.Writer) *Option {
	return &Option{
		outStream: outStream,
		errStream: errStream,
	}
}

func addCommands(rootCmd *cobra.Command, o *Option) {
	rootCmd.AddCommand(
		NewVersionCommand(o),
		NewNodesCommand(&NodesOption{}),
	)
}
