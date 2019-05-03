#!/bin/sh

set -e
echo "This is the duct-tape installation script!"

if [ ! -d /dbg ]; then
    echo "Debugging installation requires a volume mount at /dbg" 1>&2
    exit 1
fi

# Install debugging runtime files in /dbg
cd /duct-tape
cp -rp . /dbg
