![experimental](https://img.shields.io/badge/stability-experimental-orange.svg)

# Container Runtime Debugging Support Images (aka Duct Tape)

This repository defines a set of container images that package
the language runtime dependencies required to enable step-by-step
debugging of apps with
[`skaffold debug`](https://skaffold.dev/docs/how-tos/debug/). 
These container images are suitable for use as `initContainer`s on
a pod.  When executed, each container image copies these dependencies
to `/dbg/<runtimeId>`.

The idea is that `skaffold debug` will transform k8s manifests to
make available any support files required to debug specific language
runtimes.  For example, a Kubernetes podspec would be transformed to

  - create a volume to hold the debugging support files
  - run one or more of these images as `initContainer`s to populate
    this volume, mounted as `/dbg`
  - mount this volume on the applicable containers as `/dbg`
    with suitably transformed command-line in the entrypoint and arguments

Current language runtimes:

  * `go`: provides [Delve](https://github.com/go-delve/delve)
  * `python`: provides [`ptvsd`](https://github.com/Microsoft/ptvsd),
    a debug adapter that can be used for VS Code and more, for
    Python 2.7 and 3.5+
  * `nodejs`: provides a `node` wrapper that propagates `--inspect`
    args to the application invokation
  * `netcore`: provides `vsdbg` for .NET Core

## Distribution

The latest released images, which are used by `skaffold debug`, are available at:

    gcr.io/k8s-skaffold/skaffold-debug-support

Images from a particular release are available at:

    gcr.io/k8s-skaffold/skaffold-debug-support/<release>

Images from the latest commit to HEAD are available at our staging repository:

    us-central1-docker.pkg.dev/k8s-skaffold/skaffold-staging/skaffold-debug-support

You can configure Skaffold to use a specific release or the staging
repository with the following:

    skaffold config set --global debug-helpers-registry <repository>


# Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for how to contribute!
