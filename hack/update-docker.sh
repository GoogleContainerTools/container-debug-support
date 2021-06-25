#!/bin/sh
# Update Docker to 20.10.x, primarily intended for Travis
set -e

debarch() {
  case $(uname -m) in
  x86_64) echo amd64;;
  aarch64) echo arm64;;
  *) uname -m;;
  esac
}

dockerVersion=$(docker info -f '{{.ServerVersion}}')
#if [ $(echo $dockerVersion | cut -d. -f1) -lt 20 -o $(echo $dockerVersion | cut -d. -f2) -lt 10 ]; then
    distro=$(lsb_release -cs)
    arch=$(debarch)
    echo ">> Updating docker from ${dockerVersion} for $distro $arch..."
    if [ -f /etc/defaults/docker ]; then
        echo "Removing /etc/defaults/docker. Previous contents:"
        sudo sed 's/^/> /' /etc/defaults/docker
        sudo rm /etc/defaults/docker
    fi
    set -x
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=$arch] https://download.docker.com/linux/ubuntu $distro stable"
    sudo apt-get update
    sudo apt-get -y -o Dpkg::Options::="--force-confnew" install --no-install-recommends docker-ce-cli docker-ce containerd.io || (systemctl status docker.service; journalctl -xe)
#else
#    echo ">> Docker ${dockerVersion} >= 20.10"
#fi
