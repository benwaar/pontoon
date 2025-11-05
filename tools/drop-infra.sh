#!/bin/bash
# Stop and remove all containers, networks, volumes, and images defined in the infra/docker-compose.yml file

cat <<EOF
Remove all containers, networks, volumes, and images defined in the infra/docker-compose.yml file
Please type 'y' to continue or 'n' to cancel.
EOF

read -p "Are you sure you want to run the script? (y/n): " confirm
if [ "$confirm" != "y" ]; then
    echo "Script execution cancelled."
    exit 1
fi

DOCKER_COMPOSE_FILE="$(dirname "$0")/../infra/docker-compose.yml"
docker compose -f "$DOCKER_COMPOSE_FILE" down --rmi all -v