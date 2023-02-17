# This BUILD file is used to build the BCC repo.
# The BCC source code repo was added into WORKSPACE.
# Refer to this build file through the path.

load("@rules_foreign_cc//foreign_cc:defs.bzl", "cmake")

filegroup(
    name = "bcc_source",
    srcs = glob(["**"]),
)

cmake(
    name = "bcc",
    build_args = [
        "--",  # Pass options to the native tool.
        "-j`nproc`",
        "-l`nproc`",
    ],
    cache_entries = {
        "ENABLE_EXAMPLES": "OFF",
        "ENABLE_MAN": "OFF",
        "ENABLE_TESTS": "OFF",
    },
    includes = [
        "bcc/compat",
    ],
    install = False,
    lib_source = ":bcc_source",
    # These link opts are dependencies of bcc.
    linkopts = [
        # ELF binary parsing.
        "-lelf",
        # Zlib
        "-lz",
    ],
    out_static_libs = [
        "libapi-static.a",
        "libbcc.a",
        "libbcc_bpf.a",
        "libbcc-loader-static.a",
        "libclang_frontend.a",
    ],
    postfix_script = "make -C src/cc install",
    targets = [
        "api-static",
        "bcc-static",
        "bcc-loader-static",
        "bpf-static",
        "clang_frontend",
    ],
    visibility = ["//visibility:public"],
)
