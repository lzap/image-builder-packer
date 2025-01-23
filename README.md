# Image Builder Packer plugin

## Preparing the environment

Create an instance or VM with a user dedicated to image building and sudo permission to start podman (or docker) without password.

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
