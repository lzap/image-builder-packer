package ibk_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	ibk "github.com/lzap/image-builder-packer"
	"github.com/lzap/image-builder-packer/internal/sshtest"
)

func TestPodmanDryRun(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		cmd     *ibk.ContainerCliCommand
		session []sshtest.RequestReply
	}{
		{
			name: "fedora-minimal-raw-dry-run",
			cmd: &ibk.ContainerCliCommand{
				Distro:    "fedora",
				Type:      "minimal-raw",
				Arch:      "x86_64",
				Blueprint: "blueprint",
				DryRun:    true,
			},
			session: []sshtest.RequestReply{
				{
					Request: "which podman",
					Reply:   "/usr/bin/podman\n",
					Status:  0,
				},
				{
					Request: "arch",
					Reply:   "x86_64\n",
					Status:  0,
				},
				{
					Request: "mkdir ./output-hehwuXP6NyGIr",
					Reply:   "",
					Status:  0,
				},
				{
					Request: "scp -t /tmp",
					Reply:   "",
					Status:  0,
				},
				{
					Request: "echo sudo /usr/bin/podman run --privileged -v ./output-hehwuXP6NyGIr:/output -v /tmp/ibpacker-o2rHJLEEkT68y.toml:/tmp/ibpacker-o2rHJLEEkT68y.toml ghcr.io/osbuild/image-builder-cli:latest build --blueprint /tmp/ibpacker-o2rHJLEEkT68y.toml --distro fedora minimal-raw",
					Reply:   "Building...\nDone.\n",
					Status:  0,
				},
				{
					Request: "rm -f /tmp/ibpacker-o2rHJLEEkT68y.toml",
					Reply:   "",
					Status:  0,
				},
			},
		},
		{
			name: "fedora-minimal-raw",
			cmd: &ibk.ContainerCliCommand{
				Distro:    "fedora",
				Type:      "minimal-raw",
				Arch:      "x86_64",
				Blueprint: "blueprint",
			},
			session: []sshtest.RequestReply{
				{
					Request: "which podman",
					Reply:   "/usr/bin/podman\n",
					Status:  0,
				},
				{
					Request: "arch",
					Reply:   "x86_64\n",
					Status:  0,
				},
				{
					Request: "mkdir ./output-hehwuXP6NyGIr",
					Reply:   "",
					Status:  0,
				},
				{
					Request: "scp -t /tmp",
					Reply:   "",
					Status:  0,
				},
				{
					Request: "sudo /usr/bin/podman run --privileged -v ./output-hehwuXP6NyGIr:/output -v /tmp/ibpacker-o2rHJLEEkT68y.toml:/tmp/ibpacker-o2rHJLEEkT68y.toml ghcr.io/osbuild/image-builder-cli:latest build --blueprint /tmp/ibpacker-o2rHJLEEkT68y.toml --distro fedora minimal-raw",
					Reply:   "Building...\nDone.\n",
					Status:  0,
				},
				{
					Request: "rm -f /tmp/ibpacker-o2rHJLEEkT68y.toml",
					Reply:   "",
					Status:  0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibk.RandSource.Seed(0)

			server := sshtest.NewServer(t, sshtest.TestSigner(t))
			server.Handler = sshtest.RequestReplyHandler(t, tt.session)
			defer server.Close()

			buf := &ibk.CombinedWriter{}
			client, err := ibk.NewSSHTransport(ibk.SSHTransportConfig{
				Host:        server.Endpoint,
				Username:    "test",
				Password:    "unused",
				PrivateKeys: []*bytes.Buffer{bytes.NewBufferString(sshtest.PrivateKey)},
				Stdout:      buf,
				Stderr:      buf,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close(ctx)

			err = ibk.ApplyCommand(context.Background(), tt.cmd, client)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(buf.String(), "Building") {
				t.Fatalf("unexpected output: %s", buf.String())
			}
		})
	}
}
