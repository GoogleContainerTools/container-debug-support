FROM golang:1.11 as delve
RUN curl --location --output delve-1.2.0.tar.gz https://github.com/go-delve/delve/archive/v1.2.0.tar.gz \
 && tar xzf delve-1.2.0.tar.gz
RUN cd delve-1.2.0/cmd/dlv && go install

FROM python:2.7 as python27
RUN curl --location --output ptsvd.zip https://github.com/Microsoft/ptvsd/archive/v4.2.8.zip \
 && unzip -q ptsvd.zip \
 && cd ptvsd-4.2.8 \
 && mkdir -p /ptvsd/lib/python2.7/site-packages \
 && PYTHONUSERBASE=/ptvsd python setup.py install --user

FROM python:3.7 as python37
RUN curl --location --output ptsvd.zip https://github.com/Microsoft/ptvsd/archive/v4.2.8.zip \
 && unzip -q ptsvd.zip \
 && cd ptvsd-4.2.8 \
 && mkdir -p /ptvsd/lib/python3.7/site-packages \
 && PYTHONUSERBASE=/ptvsd python setup.py install --user

# Now populate the duct-tape image with the language runtime debugging support files
# The debian image is about 95MB bigger
FROM busybox
# The install script copies all files in /duct-tape to /dbg
COPY install.sh /
CMD ["/bin/sh", "/install.sh"]
WORKDIR /duct-tape
COPY --from=delve /go/bin/dlv go/bin/
COPY --from=python27 /ptvsd/ python/
COPY --from=python37 /ptvsd/ python/
