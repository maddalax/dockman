#!/bin/bash

# Pull the latest image
docker pull ghcr.io/maddalax/dockside-agent:latest

# Stop and remove the existing container, if any
docker stop dockside 2>/dev/null || true
docker rm dockside 2>/dev/null || true

# Determine the volume mount path
if [[ "$(uname)" != "Linux" ]]; then
  VOLUME_PATH="$HOME/.dockside/data"
else
  VOLUME_PATH="/data/dockside"
fi

# Run the container
docker run -d \
  --name dockside \
  --restart unless-stopped \
  -p 80:80 \
  -p 8100:3000 \
  -p 4222:4222 \
  -v "${VOLUME_PATH}:/data/dockside" \
  ghcr.io/maddalax/dockside-agent:latest
