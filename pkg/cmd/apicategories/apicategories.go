package apicategories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
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
		# Print the supported API categories & resources
		%CMD_PATH%

		# Print the supported API categories & resources sorted by a column
		%CMD_PATH% --sort-by=resource

		# Print the supported namespaced categories & resources
		%CMD_PATH% --namespaced=true

		# Print the supported non-namespaced categories & resources
		%CMD_PATH% --namespaced=false

		# Print the supported API categories & resources with a specific APIGroup
		%CMD_PATH% --api-group=extensions

		# Print the supported API categories
		%CMD_PATH% --output=category

		# Print the supported API resources in a specific API category
		%CMD_PATH% --output=resource --categories=api-extensions`)
	cmdShort = "Print the supported API resources their categories on the server"
	cmdLong  = templates.LongDesc(`
		Print the supported API resources their categories on the server.`)
)

// CmdOptions contains all the options for running the command.
type CmdOptions struct {
	Flags   *Flags
	FlagSet *pflag.FlagSet
	Client  discovery.CachedDiscoveryInterface

	genericclioptions.IOStreams
}

type sortableResourceList struct {
	list   []metav1.APIResource
	sortBy string
}

func (s sortableResourceList) Len() int      { return len(s.list) }
func (s sortableResourceList) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }
func (s sortableResourceList) Less(i, j int) bool {
	ret := strings.Compare(s.compareValues(i, j))
	if ret > 0 {
		return false
	} else if ret == 0 {
		return strings.Compare(s.list[i].Name, s.list[j].Name) < 0
	}
	return true
}

func (s sortableResourceList) compareValues(i, j int) (string, string) {
	switch s.sortBy {
	case "resource":
		return s.list[i].Name, s.list[j].Name
	default:
		return s.list[i].Group, s.list[j].Group
	}
}

// NewCmd returns an initialized Command for the command.
func NewCmd(streams genericclioptions.IOStreams, name, parentCmdPath string) *cobra.Command {
	o := &CmdOptions{
		Flags:     NewFlags(),
		IOStreams: streams,
	}

	if len(name) > 0 {
		cmdName = name
	}
	cmdPath = cmdName
	if len(parentCmdPath) > 0 {
		cmdPath = parentCmdPath + " " + cmdName
	}
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

	return cmd
}

// Complete completes all the required options for the command.
func (o *CmdOptions) Complete(cmd *cobra.Command, _ []string) error {
	var err error

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
	supportedOutputTypes := sets.NewString("", "category", "resource")
	if !supportedOutputTypes.Has(o.Flags.Output) {
		return fmt.Errorf("--output %v is not available", o.Flags.Output)
	}
	supportedSortTypes := sets.NewString("", "resource")
	if len(o.Flags.SortBy) > 0 && !supportedSortTypes.Has(o.Flags.SortBy) {
		return fmt.Errorf("--sort-by accepts only resource")
	}

	klog.V(4).Infof("Flags.APIGroup: %s", o.Flags.APIGroup)
	klog.V(4).Infof("Flags.Cached: %v", o.Flags.Cached)
	klog.V(4).Infof("Flags.Categories: %v", o.Flags.Categories)
	klog.V(4).Infof("Flags.Namespaced: %v", o.Flags.Namespaced)
	klog.V(4).Infof("Flags.NoHeaders: %v", o.Flags.NoHeaders)
	klog.V(4).Infof("Flags.Output: %s", o.Flags.Output)

	return nil
}

// Run implements all the necessary functionality for the command.
func (o *CmdOptions) Run() error {
	if _, err := o.Client.ServerVersion(); err != nil {
		return err
	}

	list, err := o.listResources()
	if err != nil {
		return err
	}

	switch o.Flags.Output {
	case "category":
		cSet := sets.NewString()
		for _, r := range list {
			for _, cat := range r.Categories {
				cSet.Insert(cat)
			}
		}
		_, err = fmt.Fprintln(o.Out, strings.Join(cSet.List(), "\n"))
	case "resource":
		rSet := sets.NewString()
		for _, r := range list {
			apiName := r.Name
			if g := r.Group; len(g) > 0 {
				apiName = fmt.Sprintf("%s.%s", r.Name, g)
			}
			rSet.Insert(apiName)
		}
		_, err = fmt.Fprintln(o.Out, strings.Join(rSet.List(), "\n"))
	default:
		var errs []error
		w := printers.GetNewTabWriter(o.Out)
		defer w.Flush()

		// print header
		if !o.Flags.NoHeaders {
			columnNames := []string{"RESOURCE", "APIGROUP", "NAMESPACED", "CATEGORIES"}
			if _, err := fmt.Fprintf(w, "%s\n", strings.Join(columnNames, "\t")); err != nil {
				errs = append(errs, err)
			}
		}

		// print rows
		sort.Stable(sortableResourceList{list: list, sortBy: o.Flags.SortBy})
		for _, r := range list {
			sortedCategories := sets.NewString(r.Categories...).List()
			if _, err := fmt.Fprintf(w, "%s\t%s\t%v\t%v\n",
				r.Name,
				r.Group,
				r.Namespaced,
				sortedCategories); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			err = errors.NewAggregate(errs)
		}
	}

	return err
}

func (o *CmdOptions) listResources() ([]metav1.APIResource, error) {
	lists, err := o.Client.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	groupChanged := o.FlagSet.Changed(flagAPIGroup)
	nsChanged := o.FlagSet.Changed(flagNamespaced)

	var resources []metav1.APIResource
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
			resources = append(resources, resource)
		}
	}

	return resources, nil
}
