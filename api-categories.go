package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
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
		%CMD_PATH%

		# Print the resources in a specific API category
		%CMD_PATH% all`)
	cmdShort = "Print the available API categories on the server"
	cmdLong  = templates.LongDesc(`
		Print the available API categories on the server.`)
)

// CmdOptions contains all the options for running the command.
type CmdOptions struct {
	Flags   *Flags
	FlagSet *pflag.FlagSet
	Client  discovery.CachedDiscoveryInterface

	RequestCategory string

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
		Args:                  cobra.MaximumNArgs(1),
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
func (o *CmdOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	//nolint:gocritic
	switch len(args) {
	case 1:
		o.RequestCategory = args[0]
	}

	// Setup flag set
	o.FlagSet = cmd.Flags()

	// Setup client
	o.Client, err = o.Flags.ToDiscoveryClient()
	if err != nil {
		return err
	}
	if !o.Flags.Cached {
		o.Client.Invalidate()
	}

	return nil
}

// Validate validates all the required options for the command.
func (o *CmdOptions) Validate() error {
	return nil
}

// Run implements all the necessary functionality for the command.
func (o *CmdOptions) Run() error {
	if _, err := o.Client.ServerVersion(); err != nil {
		return err
	}

	if len(o.RequestCategory) > 0 {
		return o.listCategoryResources(o.RequestCategory)
	}

	return o.listCategories()
}

func (o *CmdOptions) listResources() ([]*metav1.APIResource, error) {
	lists, err := o.Client.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	groupChanged := o.FlagSet.Changed(flagAPIGroup)
	nsChanged := o.FlagSet.Changed(flagNamespaced)

	var resources []*metav1.APIResource
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			klog.V(4).Infof("Ignoring invalid discovered resource %q: %v", list.GroupVersion, err)
			continue
		}
		for ix := range list.APIResources {
			resource := list.APIResources[ix]
			if len(resource.Verbs) == 0 {
				continue
			}
			if groupChanged && o.Flags.APIGroup != gv.Group {
				continue
			}
			if nsChanged && o.Flags.Namespaced != resource.Namespaced {
				continue
			}
			if len(o.Flags.Categories) > 0 && !sets.NewString(resource.Categories...).HasAll(o.Flags.Categories...) {
				continue
			}
			resource.Group = gv.Group
			resources = append(resources, &resource)
		}
	}

	return resources, nil
}

func (o *CmdOptions) listCategories() error {
	list, err := o.listResources()
	if err != nil {
		return err
	}

	cSet := sets.NewString()
	for _, r := range list {
		for _, cat := range r.Categories {
			cSet.Insert(cat)
		}
	}

	var output string
	if cSet.Len() == 0 {
		output = ("No API categories found")
	} else {
		output = strings.Join(cSet.List(), "\n")
	}
	fmt.Fprintln(o.Out, output)
	return nil
}

func (o *CmdOptions) listCategoryResources(category string) error {
	list, err := o.listResources()
	if err != nil {
		return err
	}

	rSet := sets.NewString()
	for _, r := range list {
		if sets.NewString(r.Categories...).Has(category) {
			apiName := r.Name
			if g := r.Group; len(g) > 0 {
				apiName = fmt.Sprintf("%s.%s", r.Name, g)
			}
			rSet.Insert(apiName)
		}
	}

	if rSet.Len() == 0 {
		return fmt.Errorf("the server doesn't have an API category \"%s\"", category)
	}

	fmt.Fprintln(o.Out, strings.Join(rSet.List(), "\n"))
	return nil
}
