workspace(name = "tricorder")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# https://github.com/bazelbuild/rules_go#setup
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "56d8c5a5c91e1af73eca71a6fab2ced959b67c86d12ba37feedb0a2dfea441a6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18.3")

# https://github.com/bazelbuild/bazel-gazelle#setup
http_archive(
    name = "bazel_gazelle",
    sha256 = "efbbba6ac1a4fd342d5122cbdfdb82aeb2cf2862e35022c752eaddffada7c3f3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
    ],
)

# https://github.com/bazelbuild/bazel-gazelle/issues/1233#issuecomment-1094334092
load("//:bazel/go_deps.bzl", "go_deps")

# gazelle:repository_macro bazel/go_deps.bzl%go_deps
go_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# https://github.com/bazelbuild/rules_go#protobuf-and-grpc
http_archive(
    name = "com_google_protobuf",
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# https://bazelbuild.github.io/rules_foreign_cc/0.8.0/index.html
http_archive(
    name = "rules_foreign_cc",
    # As of 2022-06-24, this is obtained from running bazel without this sha256
    # and let bazel print a sha256 in the warning.
    sha256 = "6041f1374ff32ba711564374ad8e007aef77f71561a7ce784123b9b4b88614fc",
    strip_prefix = "rules_foreign_cc-0.8.0",
    url = "https://github.com/bazelbuild/rules_foreign_cc/archive/0.8.0.tar.gz",
)

load("@rules_foreign_cc//foreign_cc:repositories.bzl", "rules_foreign_cc_dependencies")

rules_foreign_cc_dependencies()

# LLVM is built from source and installed to /opt base directory.
# It's imported into bazel's new_local_repository() by llvm_deps().
#
# Alternative way of building LLVM in bazel is copying the code in:
# https://github.com/llvm/llvm-project/blob/main/utils/bazel/WORKSPACE
# We need to choose the commit carefully from llvm repo though.
# We need to pick a commit that has the right bazel configurations.
# But not too new than what's compatible with gobpf/BCC's required BCC & llvm
# version.
#
# TODO(yzhao): Consider moving to LLVM's official bazel build after fixing the
# double free bug. One needs to bring back the llvm deps and use it to build
# bazelisk build -c dbg src/bcc:bcc_{cc,go}
# Observe the double free bug.

load("//:bazel/bcc_deps.bzl", "bcc_deps", "llvm_deps")

bcc_deps()

llvm_deps()

# Add bazel rules
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    go_image_repos = "repositories",
)

go_image_repos()

load("//:bazel/container_pulls.bzl", "container_pulls")

container_pulls()

# We need to turn on the AWS S3 Bucket ACL and then make `linux-headers.tar.gz` object can be access bu public URL
http_file(
    name = "download_linux_headers_from_s3_url",
    sha256 = "c43ff01e1e65f34714154db27070851e5a9327fa73aeb57bf018fc2290b23b60",
    urls = ["https://tricorder-dev.s3.ap-northeast-1.amazonaws.com/linux-headers.tar.gz"],
)
