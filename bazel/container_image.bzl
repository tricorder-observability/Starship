load("@bazel_skylib//rules:copy_file.bzl", _copy_file = "copy_file")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", _pkg_tar = "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", _container_image = "container_image", _container_push = "container_push")
load("@io_bazel_rules_docker//go:image.bzl", _go_image = "go_image")

# Insert manual tag to various native build rules
def _add_manual_tag(kwargs):
    if "tags" not in kwargs:
        kwargs["tags"] = ["manual"]
    else:
        kwargs["tags"].append("manual")

def _add_base_image(kwargs):
    # Bazel's default base image is `distroless`, which has old GLIBC that
    # causes crash when running api-server
    # https://github.com/bazelbuild/rules_docker/blob/master/go/image.bzl
    # https://iximiuz.com/en/posts/containers-distroless-images/
    # For now just use the official ubuntu image
    # With ubuntu base, api-server_image.tar is 121MB
    # With Bazel's default distroless base, 64MB
    # TODO(yzhao): Figure out better fixes.
    if "base" not in kwargs:
        kwargs["base"] = "@base_build_image//image"

# Insert copy_file to a companying image rule
def _add_copy_file(kwargs):
    _name = kwargs["name"]
    _tar_name = _name + ".tar"
    _src_name = ":" + _tar_name
    _out_name = "out_" + _tar_name
    _copy_file(
        name = "copy_" + _name + "_tar",
        src = _src_name,
        out = _out_name,
        tags = ["manual"],
    )

def go_image(**kwargs):
    _add_base_image(kwargs)
    _add_manual_tag(kwargs)
    _go_image(**kwargs)
    _add_copy_file(kwargs)

def container_image(**kwargs):
    _add_base_image(kwargs)
    _add_manual_tag(kwargs)
    _container_image(**kwargs)
    _add_copy_file(kwargs)

def container_push(**kwargs):
    _add_manual_tag(kwargs)
    _container_push(**kwargs)

def pkg_tar(**kwargs):
    _add_manual_tag(kwargs)
    _pkg_tar(**kwargs)

# Allow some type to use native go_iamge rule for smaller image.
native_go_image = _go_image
