![experimental](https://img.shields.io/badge/stability-experimental-orange.svg)

# Container Runtime Debugging Support (aka Duct Tape)

This repository gathers additional dependencies required to debug
particular language runtimes with [`skaffold debug`](https://skaffold.dev/docs/how-tos/debug/). 
These dependencies are packaged as a set of container images suitable
for use as `initContainer`s on a pod.  When executed, a container image
copies these dependencies to `/dbg/<runtimeId>`.

The idea is that `skaffold debug` will transform k8s manifests to
make available any support files required to debug specific language
runtimes.  For example, a Kubernetes podspec would be transformed to

  - mount a volume on `/dbg` to hold the debugging support files
  - run one or more of these `initContainer`s to populate the volume
  - mount the volume on the applicable containers as `/dbg`

Current language runtimes:
  * `go`: provides [Delve](https://github.com/go-delve/delve)
  * `python`: provides [`ptvsd`](https://github.com/Microsoft/ptvsd),
    a debug adapter that can be used for VS Code and more, for
    Python 2.7 and 3.5+
  * `nodejs`: provides a `node` wrapper that propagates `--inspect`
    args to the application invokation
  * `netcore`: provides `vsdbg` for .NET Core

## Development

This project uses Skaffold's multiple config support to allow
developing for each language runtime separately.

Each language runtime is broken out to a separate directory
with a `skaffold.yaml` for development of the `duct-tape` initContainer
image.  Each image is expected to be standalone and not require
downloading content across the network.  To add support for a new
language runtime, an image definition should download the necessary
files into the container image.  The image's entrypoint should then
copy those files into place at `/dbg/<runtime>`.  The image should
be added to the respective `skaffold.yaml` and referenced within
`test/k8s-*-installation.yaml`.

# Testing

Integration tests are found in `test/`.  These tests build and
launch applications as pods that are similar to the transformed
form produced by `skaffold debug`.  To run:

```sh
sh run-its.sh
```
