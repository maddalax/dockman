#!/bin/bash

# Pull the latest image
docker pull ghcr.io/maddalax/dockman-agent:latest

# Stop and remove the existing container, if any
docker stop dockman-agent 2>/dev/null || true
docker rm dockman-agent 2>/dev/null || true

# Determine the volume mount path
if [[ "$(uname)" != "Linux" ]]; then
  VOLUME_PATH="$HOME/.dockman/data"
else
  VOLUME_PATH="/data/dockman"
fi

# Run the container
docker run -d \
  --network host \
  --name dockman-agent \
  --restart unless-stopped \
  -v "${VOLUME_PATH}:/data/dockman" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e NATS_HOST=localhost \
  ghcr.io/maddalax/dockman-agent:latest
