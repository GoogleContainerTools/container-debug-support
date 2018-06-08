#!/bin/bash

set -e
set -x

gcloud components install docker-credential-gcr
export PATH=$PATH:/usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/bin/

# Stops any left-over containers.
docker stop $(docker container ls --quiet) || true

cd github/skaffold-java

(cd skaffold-plugins-core; ./gradlew clean build integrationTest --info --stacktrace)
(cd skaffold-maven-plugin; ./mvnw clean install -B -U -X)
(cd skaffold-gradle-plugin; ./gradlew clean build --info --stacktrace)
