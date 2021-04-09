# Copyright 2021 The Skaffold Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This Dockerfile creates a debug helper base image for Python.
# It provides installations of debugpy, ptvsd, pydevd, and pydevd-pycharm
# for Python 2.7, 3.7, 3.8, and 3.9.
#
# debugpy and ptvsd are well-structured packages installed in separate
# directories under # /dbg/python/lib/pythonX.Y/site-packages and
# that do not interfere with each other.
#
# pydevd and pydevd-pycharm install a script in .../bin and both install
# .py files directly in .../lib/pythonX.Y/site-packages.  To avoid
# interference we install pydevd and pydevd-pycharm under /dbg/python/pydevd/pythonX.Y
# and /dbg/python/pydevd-pycharm/pythonX.Y

FROM python:2.7 as python27
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy
RUN PYTHONUSERBASE=/dbgpy/pydevd/python2.7 pip install --user pydevd
RUN PYTHONUSERBASE=/dbgpy/pydevd-pycharm/python2.7 pip install --user pydevd-pycharm

FROM python:3.7 as python37
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy
RUN PYTHONUSERBASE=/dbgpy/pydevd/python3.7 pip install --user pydevd
RUN PYTHONUSERBASE=/dbgpy/pydevd-pycharm/python3.7 pip install --user pydevd-pycharm

FROM python:3.8 as python38
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy
RUN PYTHONUSERBASE=/dbgpy/pydevd/python3.8 pip install --user pydevd
RUN PYTHONUSERBASE=/dbgpy/pydevd-pycharm/python3.8 pip install --user pydevd-pycharm

FROM python:3.9 as python39
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy
RUN PYTHONUSERBASE=/dbgpy/pydevd/python3.9 pip install --user pydevd
RUN PYTHONUSERBASE=/dbgpy/pydevd-pycharm/python3.9 pip install --user pydevd-pycharm

FROM golang:1.14.1 as build
COPY launcher/ .
# Produce an as-static-as-possible wrapper binary to work on musl and glibc
RUN GOPATH="" CGO_ENABLED=0 go build -o launcher -ldflags '-s -w -extldflags "-static"' .

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM --platform=$BUILDPLATFORM busybox
ARG BUILDPLATFORM

# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=python27 /dbgpy/ python/
COPY --from=python37 /dbgpy/ python/
COPY --from=python38 /dbgpy/ python/
COPY --from=python39 /dbgpy/ python/
COPY --from=build /go/launcher python/