//go:build windows

package portfolio

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func (ExecRunner) RunBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	if name != "git" {
		return CommandResult{}, fmt.Errorf("bound_command_forbidden: only Git may use retained directory authority")
	}
	pinned, err := os.Open(directoryPath)
	if err != nil {
		return CommandResult{}, fmt.Errorf("bound_command_directory_changed: %w", err)
	}
	defer pinned.Close()
	expected, err := directory.Stat()
	if err != nil {
		return CommandResult{}, fmt.Errorf("bound_command_directory_changed: %w", err)
	}
	actual, err := pinned.Stat()
	if err != nil || !os.SameFile(expected, actual) {
		return CommandResult{}, fmt.Errorf("bound_command_directory_changed: retained identity does not match the command directory")
	}
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = directoryPath
	command.Env = sanitizedGitEnvironment(os.Environ())
	return captureCommand(command)
}
