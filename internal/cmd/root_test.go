package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/cmd"
)

func executeRoot(args ...string) (stdout string, err error) {
	root, _ := cmd.NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand_NoArgs_ShowsHelp(t *testing.T) {
	out, err := executeRoot()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "Usage") || !strings.Contains(out, "exasol-saas") {
		t.Errorf("expected help output, got: %q", out)
	}
}

func TestRootCommand_HelpFlag(t *testing.T) {
	out, err := executeRoot("--help")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "--token") {
		t.Errorf("expected --token flag in help output, got: %q", out)
	}
}

func TestRootCommand_UnknownCommand(t *testing.T) {
	_, err := executeRoot("unknowncmd")
	if err == nil {
		t.Fatal("expected non-zero exit for unknown command, got nil error")
	}
}
