#!/bin/bash

set -e

# Configuration
REMOTE_USER="root" # Replace with remote SSH username
REMOTE_HOST="fedora-server" # Replace with remote host address or IP
REMOTE_PATH="/tmp/dockman-agent.tar" # Temporary path for the tar file on the remote machine
LOCAL_IMAGE_NAME="ghcr.io/maddalax/dockman-agent:latest"
CONTAINER_NAME="dockman-agent"
DOCKER_FILE_PATH="Dockerfile-agent"

# Step 1: Build the image locally (if needed)
docker build --no-cache -t "$LOCAL_IMAGE_NAME" -f "$DOCKER_FILE_PATH" .

# Step 2: Export the image to a tar file
IMAGE_TAR="dockman-agent.tar"
docker save -o "$IMAGE_TAR" "$LOCAL_IMAGE_NAME"

# Step 3: Transfer the tar file to the remote machine
scp "$IMAGE_TAR" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"

# Step 4: Load the image on the remote machine and run the container
ssh "$REMOTE_USER@$REMOTE_HOST" << EOF
  set -e

  # Stop and remove the existing container if it exists
  docker stop "$CONTAINER_NAME" 2>/dev/null || true
  docker rm "$CONTAINER_NAME" 2>/dev/null || true

  # Load the Docker image
  docker load < "$REMOTE_PATH"

  # Remove the temporary tar file
  rm -f "$REMOTE_PATH"

  # Run the container
    docker run -d \
      --network host \
      --name dockman-agent \
      --restart unless-stopped \
      -v /data/dockman:/data/dockman \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -e NATS_HOST=100.81.35.2 \
      "$LOCAL_IMAGE_NAME"
EOF

# Step 5: Cleanup local tar file
rm -f "$IMAGE_TAR"

echo "Docker image deployed and running on $REMOTE_HOST."
