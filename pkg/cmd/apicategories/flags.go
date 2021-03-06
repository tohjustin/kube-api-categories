package apicategories

import (
	goflag "flag"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/v2"
)

const (
	flagAPIGroup        = "api-group"
	flagCached          = "cached"
	flagCategories      = "categories"
	flagNamespaced      = "namespaced"
	flagNoHeaders       = "no-headers"
	flagOutput          = "output"
	flagOutputShorthand = "o"
	flagSortBy          = "sort-by"
)

// Flags composes common configuration flag structs used in the command.
type Flags struct {
	*genericclioptions.ConfigFlags

	APIGroup   string
	Cached     bool
	Categories []string
	Namespaced bool
	NoHeaders  bool
	Output     string
	SortBy     string
}

// Copy returns a copy of Flags for mutation.
func (f *Flags) Copy() Flags {
	Flags := *f
	return Flags
}

// AddFlags receives a *pflag.FlagSet reference and binds flags related to
// configuration to it.
func (f *Flags) AddFlags(flags *pflag.FlagSet) {
	f.ConfigFlags.AddFlags(flags)

	flags.StringVar(&f.APIGroup, flagAPIGroup, f.APIGroup, "Limit to resources in the specified API group")
	flags.BoolVar(&f.Cached, flagCached, f.Cached, "If false, non-namespaced resources will be returned, otherwise returning namespaced resources by default")
	flags.StringSliceVar(&f.Categories, flagCategories, f.Categories, "Limit to resources that belong to the specified categories")
	flags.BoolVar(&f.Namespaced, flagNamespaced, f.Namespaced, "Use the cached list of resources if available")
	flags.BoolVar(&f.NoHeaders, flagNoHeaders, f.NoHeaders, "When using the default output format, don't print headers (default print headers).")
	flags.StringVarP(&f.Output, flagOutput, flagOutputShorthand, f.Output, "Output format. One of: category|resource.")
	flags.StringVar(&f.SortBy, flagSortBy, f.SortBy, "If non-empty, sort list of resources using specified field. One of: name.")

	// Hide client flags to make our help command consistent with kubectl
	_ = flags.MarkHidden("namespace")

	// Setup flags for logging.
	klogFlagSet := goflag.NewFlagSet("klog", goflag.ContinueOnError)
	klog.InitFlags(klogFlagSet)
	flags.AddGoFlagSet(klogFlagSet)

	// Logs are written to standard error instead of to files
	_ = flags.Set("logtostderr", "true")

	// Hide log flags to make our help command consistent with kubectl
	_ = flags.MarkHidden("add_dir_header")
	_ = flags.MarkHidden("alsologtostderr")
	_ = flags.MarkHidden("log_backtrace_at")
	_ = flags.MarkHidden("log_dir")
	_ = flags.MarkHidden("log_file")
	_ = flags.MarkHidden("log_file_max_size")
	_ = flags.MarkHidden("logtostderr")
	_ = flags.MarkHidden("one_output")
	_ = flags.MarkHidden("skip_headers")
	_ = flags.MarkHidden("skip_log_headers")
	_ = flags.MarkHidden("stderrthreshold")
	_ = flags.MarkHidden("v")
	_ = flags.MarkHidden("vmodule")
}

// ToClient returns a client based on the flag configuration.
func (f *Flags) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	dis, err := f.ConfigFlags.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	return dis, nil
}

// NewFlags returns flags associated with command configuration, with default
// values set.
func NewFlags() *Flags {
	return &Flags{
		ConfigFlags: genericclioptions.NewConfigFlags(true),
		APIGroup:    "",
		Cached:      false,
		Categories:  []string{},
		Namespaced:  false,
		NoHeaders:   false,
		Output:      "",
		SortBy:      "",
	}
}
