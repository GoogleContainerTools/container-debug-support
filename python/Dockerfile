FROM python:2.7 as python27
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy

FROM python:3.7 as python37
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy

FROM python:3.8 as python38
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy

FROM python:3.9 as python39
RUN PYTHONUSERBASE=/dbgpy pip install --user ptvsd debugpy

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
