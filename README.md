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

## Building without Packer

To test this library directly, do:

    go run github.com/lzap/image-builder-packer/cmd/ibpacker/ -help

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
git clone github.com/lzap/image-builder-packer
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

* Command for `bootc-image-builder`
* Packer plugin interface
* Move to `osbuild` github org
* Get the code reviewed, make a release
