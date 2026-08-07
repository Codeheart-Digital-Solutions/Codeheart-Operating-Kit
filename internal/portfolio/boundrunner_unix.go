//go:build !windows

package portfolio

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	boundGitHelperArgument = "__codeheart_internal_bound_git"
	boundGitHelperEnv      = "CODEHEART_INTERNAL_BOUND_GIT_HELPER"
)

func init() {
	if os.Getenv(boundGitHelperEnv) != "1" || len(os.Args) < 2 || os.Args[1] != boundGitHelperArgument {
		return
	}
	if err := syscall.Fchdir(3); err != nil {
		fmt.Fprintf(os.Stderr, "bound Git directory unavailable: %v\n", err)
		os.Exit(126)
	}
	gitPath, err := exec.LookPath("git")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Git unavailable: %v\n", err)
		os.Exit(127)
	}
	environment := sanitizedGitEnvironment(os.Environ())
	filtered := environment[:0]
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if key != boundGitHelperEnv {
			filtered = append(filtered, entry)
		}
	}
	if err := syscall.Exec(gitPath, append([]string{"git"}, os.Args[2:]...), filtered); err != nil {
		fmt.Fprintf(os.Stderr, "bound Git execution failed: %v\n", err)
		os.Exit(127)
	}
}

func (ExecRunner) RunBound(ctx context.Context, directory *os.File, _ string, name string, args ...string) (CommandResult, error) {
	if name != "git" {
		return CommandResult{}, fmt.Errorf("bound_command_forbidden: only Git may use retained directory authority")
	}
	executable, err := os.Executable()
	if err != nil {
		return CommandResult{}, err
	}
	command := exec.CommandContext(ctx, executable, append([]string{boundGitHelperArgument}, args...)...)
	command.ExtraFiles = []*os.File{directory}
	command.Env = append(sanitizedGitEnvironment(os.Environ()), boundGitHelperEnv+"=1")
	return captureCommand(command)
}

func (ExecRunner) StartBound(ctx context.Context, directory *os.File, _ string, name string, args ...string) (*BoundCommandStream, error) {
	if name != "git" {
		return nil, fmt.Errorf("bound_command_forbidden: only Git may use retained directory authority")
	}
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	command := exec.CommandContext(ctx, executable, append([]string{boundGitHelperArgument}, args...)...)
	command.ExtraFiles = []*os.File{directory}
	command.Env = append(sanitizedGitEnvironment(os.Environ()), boundGitHelperEnv+"=1")
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, err
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		return nil, err
	}
	return &BoundCommandStream{Stdin: stdin, Stdout: stdout, Wait: func() (CommandResult, error) {
		err := command.Wait()
		return CommandResult{Stderr: stderr.Bytes()}, err
	}}, nil
}
