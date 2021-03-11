#!/bin/sh
# Buildkit requires Docker to have experimental features enabled <https://stackoverflow.com/a/57849395>
if [ -f /etc/docker/daemon.json ]; then
    jq '.+{"experimental":"enabled"}' /etc/docker/daemon.json | sudo tee /etc/docker/daemon.json
else
    echo '{"experimental":"enabled"}' | sudo tee /etc/docker/daemon.json
fi
sudo service docker restart
