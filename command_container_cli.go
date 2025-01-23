package ibk

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"al.essio.dev/pkg/shellescape"
)

// ContainerCliCommand builds an image-builder-cli command line via podman or docker
// which builds an image from a blueprint. The blueprint is pushed to the remote
// machine via SSH and the command is executed there. The output image is saved
// to the specified directory which must be created and cleanuped up beforehand.
//
// For more information see https://github.com/osbuild/image-builder-cli
type ContainerCliCommand struct {
	// Distro is the distribution name
	Distro string

	// Type is the image type
	Type string

	// Arch is the architecture, must be set to the architecture of the remote machine
	// since cross-compilation is not supported yet.
	Arch string

	// Blueprint is the full contents of a blueprint.
	Blueprint string

	// OutputDir is the directory where the output image is saved. When unset, a new directory will
	// be created in the remote machine home directory. The caller must cleanup the directory.
	OutputDir string

	// DryRun is a flag to print the command instead of executing it. Blueprint is still pushed
	// to the remote machine and then cleaned up.
	DryRun bool

	containerCmd      string
	blueprintTempfile string
}

func (c *ContainerCliCommand) Configure(ctx context.Context, t Executor) error {
	buf := &CombinedWriter{}

	// detect container runtime
	err := t.Execute(ctx, StringCommand("which podman"), WithCombinedWriter(buf))
	if err != nil {
		return fmt.Errorf("%w podman: %w", ErrConfigure, err)
	}
	c.containerCmd = buf.FirstLine()
	slog.DebugContext(ctx, "detected podman", "path", c.containerCmd)

	if c.containerCmd == "" {
		buf.Reset()
		err := t.Execute(ctx, StringCommand("which docker"), WithCombinedWriter(buf))
		if err != nil {
			return fmt.Errorf("%w docker: %w", ErrConfigure, err)
		}
		c.containerCmd = buf.FirstLine()
		slog.DebugContext(ctx, "detected docker", "path", c.containerCmd)
	}

	if c.containerCmd == "" {
		return fmt.Errorf("%w: no container runtime found", ErrConfigure)
	}

	// detect architecture
	if c.Arch != "" {
		buf.Reset()
		err = t.Execute(ctx, StringCommand("arch"), WithCombinedWriter(buf))
		if err != nil {
			return fmt.Errorf("%w arch: %w", ErrConfigure, err)
		}
		arch := buf.FirstLine()
		slog.DebugContext(ctx, "detected architecture", "arch", arch)
		if c.Arch != arch {
			return fmt.Errorf("%w architecture mismatch: %w, output: %s", ErrConfigure, err, buf.String())
		}
	}

	// create output dir if not set
	if c.OutputDir == "" {
		buf.Reset()
		c.OutputDir = fmt.Sprintf("./output-%s", RandomString(13))
		err = t.Execute(ctx, StringCommand("mkdir "+c.OutputDir), WithCombinedWriter(buf))
		if err != nil {
			return fmt.Errorf("%w mktemp: %w, output: %s", ErrConfigure, err, buf.String())
		}
		slog.DebugContext(ctx, "created output directory", "dir", c.OutputDir)
	}

	return nil
}

func (c *ContainerCliCommand) Push(ctx context.Context, pusher Pusher) error {
	var err error

	c.blueprintTempfile, err = pusher.Push(ctx, c.Blueprint, "toml")
	return err
}

func (c *ContainerCliCommand) Build() string {
	sb := strings.Builder{}

	if c.DryRun {
		sb.WriteString("echo")
		sb.WriteRune(' ')
	}

	sb.WriteString("sudo")
	sb.WriteRune(' ')
	sb.WriteString(c.containerCmd)
	sb.WriteRune(' ')
	sb.WriteString("run --privileged")
	sb.WriteRune(' ')
	sb.WriteString("-v " + shellescape.Quote(c.OutputDir+":/output"))
	sb.WriteRune(' ')
	sb.WriteString("-v " + shellescape.Quote(c.blueprintTempfile+":"+c.blueprintTempfile))
	sb.WriteRune(' ')
	sb.WriteString("ghcr.io/osbuild/image-builder-cli:latest")
	sb.WriteRune(' ')
	sb.WriteString("build")
	sb.WriteRune(' ')
	sb.WriteString("--blueprint " + c.blueprintTempfile)
	sb.WriteRune(' ')
	sb.WriteString("--distro " + shellescape.Quote(c.Distro))
	sb.WriteRune(' ')
	sb.WriteString(shellescape.Quote(c.Type))

	return sb.String()
}
