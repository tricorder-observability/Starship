load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def bcc_deps():
    new_git_repository(
        name = "com_github_iovisor_bcc",
        # Forked from iovisor/bcc. Build instructions were updated to produced
        # additional static libraries.
        remote = "https://github.com/tricorder-observability/bcc.git",
        commit = "50de7107d6a48fcfe4f82d33433960f965d1a16a",
        shallow_since = "1675604973 +0000",
        init_submodules = True,
        recursive_init_submodules = True,
        build_file = "//:bazel/external/bcc.BUILD",
    )

def llvm_deps():
    http_archive(
        name = "llvm",
        sha256 = "f40df2ec48bfc6f3812b21b0c94cb6c23a4be47658efddf930cc6d14c2347f01",
        url = "https://tricorder-dev.s3.ap-northeast-1.amazonaws.com/starship-clang.tar.gz",
        build_file = "//:bazel/external/llvm.BUILD",
        strip_prefix = "clang-14.0",
    )
