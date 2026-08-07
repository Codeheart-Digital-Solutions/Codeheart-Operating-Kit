//go:build windows

package portfolio

import (
	"bytes"
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

func (ExecRunner) StartBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (*BoundCommandStream, error) {
	if name != "git" {
		return nil, fmt.Errorf("bound_command_forbidden: only Git may use retained directory authority")
	}
	pinned, err := os.Open(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("bound_command_directory_changed: %w", err)
	}
	expected, err := directory.Stat()
	if err != nil {
		_ = pinned.Close()
		return nil, fmt.Errorf("bound_command_directory_changed: %w", err)
	}
	actual, err := pinned.Stat()
	if err != nil || !os.SameFile(expected, actual) {
		_ = pinned.Close()
		return nil, fmt.Errorf("bound_command_directory_changed: retained identity does not match the command directory")
	}
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = directoryPath
	command.Env = sanitizedGitEnvironment(os.Environ())
	stdin, err := command.StdinPipe()
	if err != nil {
		_ = pinned.Close()
		return nil, err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		_ = pinned.Close()
		return nil, err
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		_ = pinned.Close()
		return nil, err
	}
	return &BoundCommandStream{Stdin: stdin, Stdout: stdout, Wait: func() (CommandResult, error) {
		runErr := command.Wait()
		current, statErr := os.Stat(directoryPath)
		_ = pinned.Close()
		if statErr != nil || !os.SameFile(expected, current) {
			return CommandResult{Stderr: stderr.Bytes()}, fmt.Errorf("bound_command_directory_changed: command directory identity changed")
		}
		return CommandResult{Stderr: stderr.Bytes()}, runErr
	}}, nil
}
