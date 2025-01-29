//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest mapstructure-to-hcl2 -type Config,BuildHost,AWSUpload

package main

import (
	"context"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	ibk "github.com/lzap/packer-plugin-image-builder"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// BuildHost is the host where the image will be built
	BuildHost BuildHost `mapstructure:"build_host,required"`

	// Common configuration
	ImageType    string `mapstructure:"image_type,required"`
	Architecture string `mapstructure:"architecture"`
	Blueprint    string `mapstructure:"blueprint"`

	// Regular image build configuration
	Distro string `mapstructure:"distro"`

	// Bootable container configuration
	ContainerRepository string `mapstructure:"container_repository"`

	AWSUpload AWSUpload `mapstructure:"aws_upload"`
}

type BuildHost struct {
	Hostname string `mapstructure:"hostname,required"`
	Username string `mapstructure:"username,required"`
}

type AWSUpload struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	AmiName         string `mapstructure:"ami_name"`
	S3Bucket        string `mapstructure:"s3_bucket"`
	Region          string `mapstructure:"region"`
}

type Builder struct {
	config Config
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		PluginType:  "image-builder",
		Interpolate: true,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	ui.Say("Connecting to the build host " + b.config.BuildHost.Username + "@" + b.config.BuildHost.Hostname)

	// create tail buffer
	tail := NewTailWriterThrough(1024, os.Stdout)

	// open SSH connection
	cfg := ibk.SSHTransportConfig{
		Host:     b.config.BuildHost.Hostname,
		Username: b.config.BuildHost.Username,
		Stdout:   tail,
		Stderr:   tail,
	}
	c, err := ibk.NewSSHTransport(cfg)
	if err != nil {
		return nil, err
	}
	defer c.Close(ctx)

	// configure the command
	cmd := &ibk.ContainerCliCommand{
		Distro:    b.config.Distro,
		Type:      b.config.ImageType,
		Arch:      b.config.Architecture,
		Blueprint: b.config.Blueprint,
		Common: ibk.CommonArgs{
			DryRun: true,
			TeeLog: true,
		},
	}

	// apply the command
	err = ibk.ApplyCommandPrint(ctx, cmd, c, ui.Say)
	if err != nil {
		return nil, err
	}

	// create artifact
	sa := &StringArtifact{}
	for _, line := range tail.LastLines(10) {
		sa.WriteString(line)
	}

	return sa, nil
}
