default: deploy

# Build the docker image locally
build:
    docker build -t maicare-backend:latest .

# Save the local image to a compressed tarball
package: build
    docker save maicare-backend:latest | gzip > maicare-backend.tar.gz

# Sync docker-compose configuration, environment variables, and the image tarball to the remote server
sync: package
    ssh root@167.86.75.250 "mkdir -p /root/maicare-backend"
    rsync -avz --progress maicare-backend.tar.gz root@167.86.75.250:/root/maicare-backend/
    rsync -avz docker-compose.dev.yml root@167.86.75.250:/root/maicare-backend/docker-compose.yml
    rsync -avz app.env.dev root@167.86.75.250:/root/maicare-backend/app.env
    rm maicare-backend.tar.gz

# Load the image and spin up the services on the remote server
deploy: sync
    ssh root@167.86.75.250 "cd /root/maicare-backend && docker load -i maicare-backend.tar.gz && rm maicare-backend.tar.gz && docker compose down && docker compose up -d"
