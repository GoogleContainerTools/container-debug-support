#!/bin/sh

set -e
echo "This is the duct-tape installation script!"

if [ ! -d /dbg ]; then
    echo "Debugging installation requires a volume mount at /dbg" 1>&2
    exit 1
fi

# Install Delve for Go
cp -p /duct-tape/dlv /dbg
