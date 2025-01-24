package ibk

import (
	"context"
	"errors"
	"log/slog"
)

// ErrConfigure is returned when an error occurs during configuration.
var ErrConfigure = errors.New("error while configuring")

// ErrPush is returned when an error occurs during pushing.
var ErrPush = errors.New("error while pushing")

// Command is an interface for a command that can be executed.
type Command interface {
	// Configure configures the command and prepares the environment.
	Configure(ctx context.Context, exec Executor) error

	// Push uploads one or more files to the remote host. Files are stored as temporary
	// files and are automatically cleaned up after the command is executed.
	Push(ctx context.Context, pusher Pusher) error

	// Build returns the command to execute. This command is executed in the remote environment.
	// The command should be a single string that can be passed to a shell. The command should
	// not include the shell invocation and must be protected against shell injection.
	Build() string
}

type CommonArgs struct {
	// DryRun is a flag to print the command instead of executing it. Blueprint is still pushed
	// to the remote machine and then cleaned up.
	DryRun bool

	// Interactive passes the --interactive flag to the container runtime.
	Interactive bool

	// TTY passes the --tty flag to the container runtime.
	TTY bool

	// TeeLog is a flag to tee the output of the command to a file named build.log for later use.
	TeeLog bool
}

// ApplyCommand configures, pushes, and executes a command.
func ApplyCommand(ctx context.Context, c Command, t Transport) error {
	err := c.Configure(ctx, t)
	if err != nil {
		return err
	}

	err = c.Push(ctx, t)
	if err != nil {
		return err
	}

	err = t.Execute(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

var ErrNoContainerRuntime = errors.New("no container runtime found")

func which(ctx context.Context, exec Executor, name ...string) (string, error) {
	buf := &CombinedWriter{}
	for _, n := range name {
		cmd := "which " + n
		slog.DebugContext(ctx, "running command", "cmd", cmd)
		err := exec.Execute(ctx, StringCommand(cmd), WithCombinedWriter(buf))
		if err != nil {
			return "", err
		}
		if binary := buf.FirstLine(); binary != "" {
			slog.InfoContext(ctx, "found executable", "exec", binary)
			return binary, nil
		}
	}
	return "", ErrNoContainerRuntime
}

func exec(ctx context.Context, exec Executor, cmd string) (string, error) {
	buf := &CombinedWriter{}
	slog.DebugContext(ctx, "running command", "cmd", cmd)
	err := exec.Execute(ctx, StringCommand(cmd), WithCombinedWriter(buf))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func tail1(ctx context.Context, exec Executor, cmd string) (string, error) {
	buf := &CombinedWriter{}
	slog.DebugContext(ctx, "running command", "cmd", cmd)
	err := exec.Execute(ctx, StringCommand(cmd), WithCombinedWriter(buf))
	if err != nil {
		return "", err
	}
	return buf.FirstLine(), nil
}
