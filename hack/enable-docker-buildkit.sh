#!/bin/sh
# https://www.docker.com/blog/multi-arch-build-what-about-travis/

echo ">> enabling experimental mode"
if [ -f /etc/docker/daemon.json ]; then
    echo "/etc/docker/daemon.json was:"
    sed 's/^/> /' /etc/docker/daemon.json
    echo "/etc/docker/daemon.json now:"
    jq '.+{"experimental":true}' /etc/docker/daemon.json | sudo tee /etc/docker/daemon.json
else
    sudo mkdir -vp /etc/docker
    echo "/etc/docker/daemon.json now:"
    echo '{"experimental":true}' | sudo tee /etc/docker/daemon.json
fi

if [ -f $HOME/.docker/config.json ]; then
    echo "$HOME/.docker/config.json was:"
    sed 's/^/> /' $HOME/.docker/config.json
    echo "$HOME/.docker/config.json now:"
    jq '.+{"experimental":"enabled"}' /etc/docker/daemon.json | tee $HOME/.docker/config.json
else 
    mkdir -vp $HOME/.docker
    echo "$HOME/.docker/config.json now:"
    echo '{"experimental":"enabled"}' | tee $HOME/.docker/config.json
fi

echo ">> updating docker engine"
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get -y -o Dpkg::Options::="--force-confnew" install docker-ce

echo ">> installing docker-buildx"
mkdir -vp $HOME/.docker/cli-plugins/
curl --silent -L "https://github.com/docker/buildx/releases/download/v0.5.1/buildx-v0.5.1.linux-amd64" > $HOME/.docker/cli-plugins/docker-buildx
chmod a+x $HOME/.docker/cli-plugins/docker-buildx

docker info
