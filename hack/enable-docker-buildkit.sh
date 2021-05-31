#!/bin/sh
# https://www.docker.com/blog/multi-arch-build-what-about-travis/

echo ">> enabling experimental mode"
if [ -f /etc/docker/daemon.json ]; then
    echo "/etc/docker/daemon.json was:"
    sed 's/^/> /' /etc/docker/daemon.json
    echo "/etc/docker/daemon.json now:"
    jq '.+{"experimental":true}' /etc/docker/daemon.json \
    | jq '."registry-mirrors" -= ["https://registry.docker.io"]' \
    | jq '."registry-mirrors" += ["https://mirror.gcr.io"]' \
    | sudo tee /etc/docker/daemon.json
else
    sudo mkdir -vp /etc/docker
    echo "/etc/docker/daemon.json now:"
    echo '{"experimental":true,"registry-mirrors":["https://mirror.gcr.io"]}' \
    | sudo tee /etc/docker/daemon.json
fi

if [ -f $HOME/.docker/config.json ]; then
    echo "$HOME/.docker/config.json was:"
    sed 's/^/> /' $HOME/.docker/config.json
    echo "$HOME/.docker/config.json now:"
    jq '.+{"experimental":"enabled"}' /etc/docker/daemon.json \
    | tee $HOME/.docker/config.json
else 
    mkdir -vp $HOME/.docker
    echo "$HOME/.docker/config.json now:"
    echo '{"experimental":"enabled"}' \
    | tee $HOME/.docker/config.json
fi

echo ">> installing docker-buildx"
mkdir -vp $HOME/.docker/cli-plugins/
curl --silent -L "https://github.com/docker/buildx/releases/download/v0.5.1/buildx-v0.5.1.linux-${TRAVIS_CPU_ARCH}" > $HOME/.docker/cli-plugins/docker-buildx
chmod a+x $HOME/.docker/cli-plugins/docker-buildx
