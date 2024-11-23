#!/bin/bash

# Pull the latest image
docker pull ghcr.io/maddalax/dockside-agent:latest

# Stop and remove the existing container, if any
docker stop dockside-agent 2>/dev/null || true
docker rm dockside-agent 2>/dev/null || true

# Determine the volume mount path
if [[ "$(uname)" != "Linux" ]]; then
  VOLUME_PATH="$HOME/.dockside/data"
else
  VOLUME_PATH="/data/dockside"
fi

# Run the container
docker run -d \
  --network host \
  --name dockside-agent \
  --restart unless-stopped \
  -v "${VOLUME_PATH}:/data/dockside" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e NATS_HOST=localhost \
  ghcr.io/maddalax/dockside-agent:latest
