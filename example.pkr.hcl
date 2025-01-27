source "image-builder" "example" {
    build_host {
        hostname = "example.com"
        username = "builder"
    }

    image_type = "raw"
    container_repository = "quay.io/centos-bootc/centos-bootc:stream9"

    blueprint = <<EOV
[[customizations.user]]
name = "user"
password = "changeme"
groups = ["wheel"]
EOV
}

build {
    sources = [ "source.image-builder.example" ]
}
