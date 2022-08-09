#!/bin/sh

docker run --rm --add-host host.docker.internal:host-gateway --publish 4566:4566 --publish 4510-4559:4510-4559 localstack/localstack:1.0.0
