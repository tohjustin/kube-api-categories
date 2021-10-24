package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var rootCmdName = "kube-api-categories"

//nolint:gochecknoinits
func init() {
	// If executed as a kubectl plugin
	if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		rootCmdName = "kubectl api-categories"
	}
}

func main() {
	flags := pflag.NewFlagSet("kube-api-categories", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	rootCmd := NewCmd(streams, rootCmdName)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
