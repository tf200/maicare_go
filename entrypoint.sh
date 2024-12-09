#!/bin/sh

# Generate app.env from environment variables
cat <<EOF > /app/app.env
DB_SOURCE=${DB_SOURCE}
SERVER_ADDRESS=${SERVER_ADDRESS}
SECRET_KEY=${SECRET_KEY}
ACCESS_TOKEN_DURATION=${ACCESS_TOKEN_DURATION}
REFRESH_TOKEN_DURATION=${REFRESH_TOKEN_DURATION}
EOF

# Optionally, you can echo the contents for debugging (remove in production)
# cat /app/app.env

# Execute the main application
exec /app/main