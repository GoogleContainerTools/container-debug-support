FROM mcr.microsoft.com/dotnet/runtime:6.0 as netcore
RUN apt-get update \
    && apt-get install -y --no-install-recommends unzip curl \
    && curl -sSL https://aka.ms/getvsdbgsh | bash /dev/stdin -v latest -l /vsdbg

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM --platform=$BUILDPLATFORM busybox
ARG BUILDPLATFORM

# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=netcore /vsdbg/ netcore/
