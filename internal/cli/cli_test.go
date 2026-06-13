package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"--version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if stdout.String() != "Dotbot-Go version 0.2.1\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestExecuteAppExitDoesNotDuplicateToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"--force-color", "--no-color"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
	if !strings.Contains(stdout.String(), "`--force-color` and `--no-color` cannot both be provided") {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestExecuteHelpUsesStyledSections(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	got := stdout.String()
	for _, expected := range []string{
		"dotbot-go",
		"Usage",
		"Examples",
		"Built-In Directives",
		"Flags",
		"Output",
		"--config-file <file>",
		"--dry-run",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("missing %q in help:\n%s", expected, got)
		}
	}
	if strings.Contains(got, "\033[") {
		t.Fatalf("unexpected color for buffer output: %q", got)
	}
	if strings.Contains(got, "Compatibility") || strings.Contains(got, "--plugin") {
		t.Fatalf("help includes plugin support: %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestExecuteRejectsPluginFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"--plugin", "example.py"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "unknown flag: --plugin") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestExecuteHelpForceColor(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"--force-color", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "\033[1;36mdotbot-go\033[0m") {
		t.Fatalf("missing colored title: %q", stdout.String())
	}
}
