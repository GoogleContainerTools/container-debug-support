#!/bin/sh
# Update Docker to 20.10.x for Ubuntu (Travis)
set -e

debarch() {
  case $(uname -m) in
  x86_64) echo amd64;;
  aarch64) echo arm64;;
  *) uname -m;;
  esac
}

dockerVersion=$(docker info -f '{{.ServerVersion}}')
echo "Checking docker server version: ${dockerVersion}"
if [ $(echo $dockerVersion | cut -d. -f1) -lt 20 -o $(echo $dockerVersion | cut -d. -f2) -lt 10 ]; then

    echo ">> Updating docker for host arch $(debarch)..."
    set -x
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=$(debarch)] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get -y -o Dpkg::Options::="--force-confnew" install docker-ce-cli docker-ce containerd.io
fi
