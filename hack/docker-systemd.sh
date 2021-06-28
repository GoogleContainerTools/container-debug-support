#!/bin/sh

# change docker runtime to use systemd https://kubernetes.io/docs/setup/production-environment/container-runtimes/
cat <<EOF >/tmp/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

sudo cp /tmp/daemon.json /etc/docker/daemon.json
