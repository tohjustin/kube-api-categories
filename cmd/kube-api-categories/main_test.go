package main_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	kubeapicategories "github.com/tohjustin/kube-api-categories/cmd/kube-api-categories"
	"github.com/tohjustin/kube-api-categories/internal/version"
)

func runCmd(args ...string) (string, error) {
	buf := bytes.NewBufferString("")
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	cmd := kubeapicategories.NewCmd(streams)
	cmd.SetOut(buf)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		return "", err
	}
	out, err := io.ReadAll(buf)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func TestCommandWithVersionFlag(t *testing.T) {
	t.Parallel()

	output, err := runCmd("--version")
	if err != nil {
		t.Fatalf("failed to run command: %v", err)
	}

	expected := fmt.Sprintf("%#v\n", version.Get())
	if output != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, output)
	}
}
