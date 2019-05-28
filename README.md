![experimental](https://img.shields.io/badge/stability-experimental-orange.svg)

# Container Runtime Debugging Support (aka Duct Tape)

_Caution: this is work-in-progress_

This repository gathers additional dependencies required to debug
particular language runtimes with [`skaffold debug`](https://skaffold.dev/docs/how-tos/debug/). 
These dependencies are packaged as a set of container images suitable
for use as `initContainer`s on a pod.

Currently supported:
  * Go: provides [Delve](https://github.com/go-delve/delve)
  * Python: provides [`ptvsd`](https://github.com/Microsoft/ptvsd),
    a debug adapter that can be used for VS Code and more

## Development

The idea is that `skaffold debug` will transform k8s manifests to
make available any support files required to debug specific language
runtimes.  For example, the `k8s-pod.yaml.orig` file shows an k8s
podspec that would be transformed to `k8s-pod.yaml` to:
  - mount a volume to hold the debugging support files
  - provide an initContainer to populate the volume
  - mount the volume to the applicable containers

This directory includes a `skaffold.yaml` for development of the
the `duct-tape` initContainer image.  To add support for a new
language runtime, run `skaffold dev` and tweak `duct-tape/Dockerfile`
to download and install the necessary files in `/duct-tape`.  The
initContainer will then copy the contents of this image into place
via its entrypoint (`install.sh`).
