FROM --platform=$BUILDPLATFORM curlimages/curl as netcore
ARG BUILDPLATFORM
ARG TARGETPLATFORM
# assume glibc; RuntimeIDs gleaned from the getvsdbgsh script
RUN RuntimeID=$(case "$TARGETPLATFORM" in linux/amd64) echo linux-x64;; linux/arm64) echo linux-arm64;; *) exit 1;; esac); \
 mkdir $HOME/vsdbg && curl -sSL https://aka.ms/getvsdbgshbeta | sh /dev/stdin -v latest -l $HOME/vsdbg -r $RuntimeID

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM --platform=$TARGETPLATFORM busybox
ARG TARGETPLATFORM

# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=netcore /home/curl_user/vsdbg/ netcore/
