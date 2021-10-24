package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	cmdPath    string
	cmdName    = "api-categories"
	cmdUse     = "%CMD% [options]"
	cmdExample = templates.Examples(`
		# Print the available API categories
		%CMD_PATH%`)
	cmdShort = "Print the available API categories on the server"
	cmdLong  = templates.LongDesc(`
		Print the available API categories on the server.`)
)

// CmdOptions contains all the options for running the command.
type CmdOptions struct {
	Flags  *Flags
	Client discovery.DiscoveryInterface

	genericclioptions.IOStreams
}

// NewCmd returns an initialized Command for the command.
func NewCmd(streams genericclioptions.IOStreams, name string) *cobra.Command {
	o := &CmdOptions{
		Flags:     NewFlags(),
		IOStreams: streams,
	}

	if len(name) > 0 {
		cmdName = name
	}
	cmdPath = cmdName
	cmd := &cobra.Command{
		Use:                   strings.ReplaceAll(cmdUse, "%CMD%", cmdName),
		Example:               strings.ReplaceAll(cmdExample, "%CMD_PATH%", cmdPath),
		Short:                 cmdShort,
		Long:                  cmdLong,
		Args:                  cobra.MaximumNArgs(0),
		DisableFlagsInUseLine: true,
		DisableSuggestions:    true,
		SilenceUsage:          true,
		Run: func(c *cobra.Command, args []string) {
			klog.V(4).Infof("Version: %s", c.Root().Version)
			cmdutil.CheckErr(o.Complete(c, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	// Setup flags
	o.Flags.AddFlags(cmd.Flags())

	// Setup version flag
	cmd.SetVersionTemplate("{{printf \"%s\" .Version}}\n")
	cmd.Version = fmt.Sprintf("%#v", GetVersion())

	return cmd
}

// Complete completes all the required options for the command.
func (o *CmdOptions) Complete(_ *cobra.Command, _ []string) error {
	var err error

	// Setup client
	o.Client, err = o.Flags.ToDiscoveryClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate validates all the required options for the command.
func (o *CmdOptions) Validate() error {
	return nil
}

// Run implements all the necessary functionality for the command.
func (o *CmdOptions) Run() error {
	// First check if Kubernetes cluster is reachable
	if _, err := o.Client.ServerVersion(); err != nil {
		return err
	}

	return nil
}
