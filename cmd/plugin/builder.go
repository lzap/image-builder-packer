//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest mapstructure-to-hcl2 -type Config,BuildHost,AWSUpload

package main

import (
	"context"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	BuildHost BuildHost `mapstructure:"build_host,required"`

	ImageType    string `mapstructure:"image_type,required"`
	Architecture string `mapstructure:"architecture"`
	Blueprint    string `mapstructure:"blueprint"`

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
	ui.Say("Started building image")

	return &artifact{}, nil
}
