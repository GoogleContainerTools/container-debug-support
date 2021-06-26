#!/bin/sh

echo ">> enabling gcr mirror"
if [ -f /etc/docker/daemon.json ]; then
    echo "/etc/docker/daemon.json was:"
    sed 's/^/> /' /etc/docker/daemon.json

    # avoid `jq ... | sudo tee /etc/docker/daemon.json` as we were
    # having 0-byte files created instead (!)
    cat /etc/docker/daemon.json \
    | jq '."registry-mirrors" -= ["https://registry.docker.io"]' \
    | jq '."registry-mirrors" += ["https://mirror.gcr.io"]' \
    | tee /tmp/docker-daemon.json
    sudo cp /tmp/docker-daemon.json /etc/docker/daemon.json

    echo "/etc/docker/daemon.json now:"
    sed 's/^/> /' /etc/docker/daemon.json
else
    sudo mkdir -vp /etc/docker
    echo "/etc/docker/daemon.json now:"
    echo '{"registry-mirrors":["https://mirror.gcr.io"]}' \
    | sudo tee /etc/docker/daemon.json
fi

