load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

def _docker_hub(**kwargs):
    kwargs["registry"] = "index.docker.io"
    container_pull(**kwargs)

def container_pulls():
    _docker_hub(
        name = "official_ubuntu",
        repository = "library/ubuntu",
        tag = "22.04",
        digest = "sha256:965fbcae990b0467ed5657caceaec165018ef44a4d2d46c7cdea80a9dff0d1ea",
    )

    _docker_hub(
        name = "postgres",
        repository = "postgres",
        tag = "15",
        digest = "sha256:bc068d3dfeb1185159a22a325a788b749a681285d428b206f9fa1a73c29b4dd8",
    )

    _docker_hub(
        name = "promscale",
        repository = "timescale/promscale",
        tag = "15",
        digest = "sha256:c249987e225d3c6570134b7f64233e0b70bb4e28ab1ef6f245616388d9510607",
    )

    # This image is too big, the .tar file is 2.4 GB, so we just use the vanilla
    # timescaledb.
    # _docker_hub(
    #     name = "timescaledb-ha",
    #     repository = "timescale/timescaledb-ha",
    #     tag = "pg14-latest",
    #     digest = "sha256:c6879ffd8d6167c82b2aa8df1d3941f586ce37ad485935cb6a46b68dd46b6b53",
    # )

    _docker_hub(
        name = "timescaledb",
        repository = "timescale/timescaledb",
        tag = "2.9.2-pg14",
        digest = "sha256:b124adba582ea65174cbbe4ab86fa4ad7125b10442cdcc5bb87f739bce6ce35b",
    )

    # Base image used by bazel to build container images.
    # See $ToT/devops/docker/base_build_image/README.md for more details.
    container_pull(
        name = "base_build_image",
        registry = "public.ecr.aws/tricorder",
        repository = "base_build_image",
        tag = "v0.1",
        digest = "sha256:04bc4b96ca3763ad37d04b4914aaede13ff47ef1d5c6e6d7f17dba6da4d5f421",
    )
