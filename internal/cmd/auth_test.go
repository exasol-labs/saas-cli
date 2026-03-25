package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/exasol-labs/saas-cli/internal/cmd"
	"github.com/spf13/cobra"
)

func executeWithAPICmd(args ...string) (output string, cfg *cmd.Config, err error) {
	root, cfg := cmd.NewRootCmd()

	dummy := &cobra.Command{
		Use:   "dummy",
		Short: "Dummy API subcommand for testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	root.AddCommand(dummy)

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), cfg, err
}

func TestAuth_TokenFromEnvVar(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "env-token-value")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "test-account-id")

	_, _, err := executeWithAPICmd("dummy")
	if err != nil {
		t.Fatalf("expected no error when EXASOL_SAAS_TOKEN is set, got: %v", err)
	}
}

func TestAuth_TokenFromFlag(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "test-account-id")

	_, _, err := executeWithAPICmd("--token", "flag-token-value", "dummy")
	if err != nil {
		t.Fatalf("expected no error when --token flag is set, got: %v", err)
	}
}

func TestAuth_FlagPrecedenceOverEnvVar(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "env-token")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "test-account-id")

	_, cfg, err := executeWithAPICmd("--token", "flag-token", "dummy")
	if err != nil {
		t.Fatalf("expected no error when both --token and env var are set, got: %v", err)
	}
	if cfg.Token != "flag-token" {
		t.Errorf("expected cfg.Token to be %q (flag wins), got %q", "flag-token", cfg.Token)
	}
}

func TestAuth_NoToken_ExitsWithError(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "test-account-id")

	out, _, err := executeWithAPICmd("dummy")
	if err == nil {
		t.Fatal("expected error when no token is provided, got nil")
	}
	if !strings.Contains(out+err.Error(), "token") {
		t.Errorf("expected error message to mention token, got output: %q, err: %v", out, err)
	}
}

func TestAuth_AccountIDFromEnvVar(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "test-token")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "env-account-id")

	_, cfg, err := executeWithAPICmd("dummy")
	if err != nil {
		t.Fatalf("expected no error when EXASOL_SAAS_ACCOUNT_ID is set, got: %v", err)
	}
	if cfg.AccountID != "env-account-id" {
		t.Errorf("expected cfg.AccountID to be %q, got %q", "env-account-id", cfg.AccountID)
	}
}

func TestAuth_AccountIDFromFlag(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "test-token")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "")

	_, cfg, err := executeWithAPICmd("--account-id", "flag-account-id", "dummy")
	if err != nil {
		t.Fatalf("expected no error when --account-id flag is set, got: %v", err)
	}
	if cfg.AccountID != "flag-account-id" {
		t.Errorf("expected cfg.AccountID to be %q, got %q", "flag-account-id", cfg.AccountID)
	}
}

func TestAuth_AccountIDFlagPrecedenceOverEnvVar(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "test-token")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "env-account-id")

	_, cfg, err := executeWithAPICmd("--account-id", "flag-account-id", "dummy")
	if err != nil {
		t.Fatalf("expected no error when both --account-id and env var are set, got: %v", err)
	}
	if cfg.AccountID != "flag-account-id" {
		t.Errorf("expected cfg.AccountID to be %q (flag wins), got %q", "flag-account-id", cfg.AccountID)
	}
}

func TestAuth_NoAccountID_ExitsWithError(t *testing.T) {
	t.Setenv("EXASOL_SAAS_TOKEN", "test-token")
	t.Setenv("EXASOL_SAAS_ACCOUNT_ID", "")

	out, _, err := executeWithAPICmd("dummy")
	if err == nil {
		t.Fatal("expected error when no account ID is provided, got nil")
	}
	if !strings.Contains(out+err.Error(), "account") {
		t.Errorf("expected error message to mention account, got output: %q, err: %v", out, err)
	}
}
