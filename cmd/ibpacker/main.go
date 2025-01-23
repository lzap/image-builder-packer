// Copyright 2025 Red Hat Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	ibk "github.com/osbuild/image-builder-packer"
)

// A utility for testing the code without packer integration

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100000000*time.Second)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	var hostname = flag.String("hostname", "", "SSH hostname")
	var port = flag.Int("port", 22, "SSH port")
	var username = flag.String("username", "", "SSH username")
	var distro = flag.String("distro", "fedora", "distribution name (fedora, centos, rhel, ...)")
	var imageType = flag.String("type", "minimal-raw", "image type (minimal-raw, qcow2, ...)")
	var arch = flag.String("arch", "x86_64", "architecture")
	var blueprintFile = flag.String("blueprint", "", "path to blueprint file")
	var dryRun = flag.Bool("dry-run", false, "dry run")
	flag.Parse()

	cfg := ibk.SSHTransportConfig{
		Host:     *hostname,
		Port:     *port,
		Username: *username,
	}

	c, err := ibk.NewSSHTransport(cfg)
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)

	// load blueprint into a string
	blueprint, err := os.ReadFile(*blueprintFile)
	if err != nil {
		panic(err)
	}

	cmd := &ibk.ContainerCliCommand{
		Distro:    *distro,
		Type:      *imageType,
		Arch:      *arch,
		Blueprint: string(blueprint),
		DryRun:    *dryRun,
	}
	err = ibk.ApplyCommand(ctx, cmd, c)
	if err != nil {
		panic(err)
	}
}
