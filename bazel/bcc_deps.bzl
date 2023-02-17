load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

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
    native.new_local_repository(
        name = "llvm",
        build_file = "bazel/external/llvm.BUILD",
        path = "/opt/tricorder/lib/clang-14.0",
    )
