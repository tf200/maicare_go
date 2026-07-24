default: deploy

remote_db_url := "postgresql://maicare:maicare@167.86.75.250:25432/maicare?sslmode=disable"

# Build the docker image locally
build:
    docker build -t maicare-backend:latest .

# Save the local image to a compressed tarball
package: build
    docker save maicare-backend:latest | gzip > maicare-backend.tar.gz

# Sync docker-compose configuration, environment variables, migrations, and the image tarball to the remote server
sync: package
    ssh root@167.86.75.250 "mkdir -p /root/maicare-backend/db"
    rsync -avz --progress maicare-backend.tar.gz root@167.86.75.250:/root/maicare-backend/
    rsync -avz docker-compose.dev.yml root@167.86.75.250:/root/maicare-backend/docker-compose.yml
    rsync -avz app.env.dev root@167.86.75.250:/root/maicare-backend/app.env
    rsync -avz db/migrations root@167.86.75.250:/root/maicare-backend/db/
    rm maicare-backend.tar.gz

# Load the image and spin up the services on the remote server
deploy: sync
    ssh root@167.86.75.250 "cd /root/maicare-backend && docker load -i maicare-backend.tar.gz && rm maicare-backend.tar.gz && docker compose down && docker compose up -d"

# Run database migration up on the remote server (default 1 step)
migrate-up count="1":
    migrate -path db/migrations -database "{{remote_db_url}}" -verbose up {{count}}

# Run database migration down on the remote server (default 1 step)
migrate-down count="1":
    migrate -path db/migrations -database "{{remote_db_url}}" -verbose down {{count}}

# Force database migration version on the remote server
migrate-force version:
    migrate -path db/migrations -database "{{remote_db_url}}" force {{version}}

# Run database migration, roles sync, seed data, and create admin user on the remote server
seed-all: migrate-up roles-sync seed create-admin

# Sync roles and permissions on the remote server
roles-sync:
    ssh root@167.86.75.250 "docker compose -f /root/maicare-backend/docker-compose.yml exec -T app ./roles-sync"

# Seed the remote database with mock data
seed:
    ssh root@167.86.75.250 "docker compose -f /root/maicare-backend/docker-compose.yml exec -T app ./seed"

# Create the admin user on the remote server
create-admin:
    ssh root@167.86.75.250 "docker compose -f /root/maicare-backend/docker-compose.yml exec -T app ./create-admin"

