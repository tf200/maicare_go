#!/bin/sh

# Generate app.env from environment variables
# cat <<EOF > /app/app.env
# DB_SOURCE=${DB_SOURCE}
# SERVER_ADDRESS=${SERVER_ADDRESS}
# ACCESS_TOKEN_SECRET_KEY=${ACCESS_TOKEN_SECRET_KEY}
# REFRESH_TOKEN_SECRET_KEY=${REFRESH_TOKEN_SECRET_KEY}
# ACCESS_TOKEN_DURATION=${ACCESS_TOKEN_DURATION}
# REFRESH_TOKEN_DURATION=${REFRESH_TOKEN_DURATION}
# B2_KEY=${B2_KEY}
# B2_KEY_ID=${B2_KEY_ID}
# B2_BUCKET=${B2_BUCKET}
# HOST=${HOST}
# EOF

# # Echo the contents for debugging
# echo "Generated app.env contents:"
# cat /app/app.env

# Execute the main application
exec /app/main
