#!/bin/sh

read stdin
echo $stdin

echo "output"

>&2 echo "error"
