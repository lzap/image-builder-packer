# Image Builder Packer plugin

HashiCorp [Packer](https://www.packer.io/) plugin for [image-builder-cli](https://github.com/osbuild/image-builder-cli) and [bootc-image-builder](https://github.com/osbuild/bootc-image-builder).

## Preparing the environment

All building is done remotely over SSH connection, you need to create an instance or VM with a user dedicated to image building and sudo permission to start podman (or docker) without password.

    adduser -m builder

Either setup a password

    passwd builder

or preferably deploy a public SSH key (execute from machine with packer)

    ssh-copy-id builder@host

Make sure the container runtime can be executed without password.

```
cat <<EOF >/etc/sudoers.d/builder
builder ALL=(ALL) NOPASSWD: /usr/bin/podman, /usr/bin/docker
EOF
```

Cross-architecture building is currently not supported.

## Building using image-builder-cli

Install packer *on your machine* not on the builder instance/VM, for example on Fedora:

```
sudo dnf install -y dnf-plugins-core
sudo dnf config-manager addrepo --from-repofile=https://rpm.releases.hashicorp.com/fedora/hashicorp.repo
sudo dnf -y install packer
```

Create a packer template named `template.pkr.hcl`:

```
packer {
  required_plugins {
    image-builder = {
      source = "github.com/lzap/image-builder"
      version = ">= 0.0.1"
    }
  }
}

source "image-builder" "example" {
    build_host {
        hostname = "zzzap.tpb.lab.eng.brq.redhat.com"
        username = "builder"
    }

    container_repository = "quay.io/centos-bootc/centos-bootc:stream9"

    blueprint = <<EOV
[[customizations.user]]
name = "user"
password = "changeme"
groups = ["wheel"]
EOV

    image_type = "raw"
}

build {
    sources = [ "source.image-builder.example" ]
}
```

Perform the build via:

      packer init template.pkr.hcl
      packer build template.pkr.hcl

The image builder plugin will print last several lines from the image builder output as an artifact. To see more detailed output:

      PACKER_LOG=1 packer build template.pkr.hcl

## Dry run

If you want to perform, for any reason, a dry run where the main build command is `echo`ed to the console rather than executed, just set `IMAGE_BUILDER_DRY_RUN=true` environment variable when executing packer.

## Building without Packer

To test this library directly without packer, do:

    go run github.com/lzap/packer-plugin-image-builder/cmd/ibpacker/ -help

Use options to initiate a build:

```
Usage of ibpacker:
  -arch string
        architecture (default "x86_64")
  -blueprint string
        path to blueprint file
  -distro string
        distribution name (fedora, centos, rhel, ...) (default "fedora")
  -dry-run
        dry run
  -hostname string
        SSH hostname or IP with optional port (e.g. example.com:22)
  -type string
        image type (minimal-raw, qcow2, ...) (default "minimal-raw")
  -username string
        SSH username
```

For example:

```
git clone github.com/lzap/packer-plugin-image-builder
go run ./cmd/ibpacker/ \
    -hostname example.com \
    -username builder \
    -distro centos-9 \
    -type minimal-raw \
    -blueprint ./cmd/ibpacker/blueprint_example.toml
```

See [osbuild blueprint reference](https://osbuild.org/docs/user-guide/blueprint-reference/) for more info about blueprint format.

## Testing

To run unit and integration test against mock SSH server running on localhost:

    go test .

Keys in `internal/sshtest/keys.go` are just dummy (test only) keys, you may receive false positives from security scanners about leaked keys when cloning the repo.

## LICENSE

Apache Version 2.0

## TODO

* Integration test via packer command with SSH mock server
* Move to `osbuild` github org
* Get the code reviewed, make a release
