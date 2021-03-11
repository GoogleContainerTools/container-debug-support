#!/bin/sh
# Buildkit requires Docker to have experimental features enabled <https://stackoverflow.com/a/57849395>

if [ -f /etc/docker/daemon.json ]; then
    echo "/etc/docker/daemon.json was:"
    sed 's/^/> /' /etc/docker/daemon.json
    echo "/etc/docker/daemon.json now:"
    jq '.+{"experimental":true}' /etc/docker/daemon.json | sudo tee /etc/docker/daemon.json
else
    echo "/etc/docker/daemon.json now:"
    echo '{"experimental":true}' | sudo tee /etc/docker/daemon.json
fi
sudo chown travis:travis /etc/docker/daemon.json
echo "Restarting docker..."
sudo systemctl restart docker || (echo "Failed!"; sudo journalctl -xe; exit 1)
sudo systemctl status docker.service

docker info
