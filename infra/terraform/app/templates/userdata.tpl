#!/bin/bash
cat > /home/ec2-user/serve/server/.env <<EOF
DEV_MODE=${dev_mode}
SERVE_DAY=${serve_day}

# Server config
PORT=${api_port}

# Database config
PGHOST=${db_host}
PGPORT=${db_port}
PGUSER=${db_user}
PGPASSWORD=${db_pass}
PGDATABASE=${db_name}
DATABASE_URL=

# Auth0 config
AUTH0_DOMAIN=${auth0_domain}
AUTH0_AUDIENCE=${auth0_audience}
AUTH0_CLIENT_ID=${auth0_client_id}
AUTH0_CLIENT_SECRET=${auth0_client_secret}

# MailTrap config
MAIL_HOST=${mailtrap_host}
MAIL_KEY=${mailtrap_key}
MAIL_FROM=${mailtrap_from}
MAIL_PORT=${mailtrap_port}
MAIL_USER=${mailtrap_user}
MAIL_PASS=${mailtrap_pass}

# Clearstream Text config
CS_API_KEY=${clearstream_api_key}
CS_TEXT_FROM=${clearstream_text_from}
# Google Maps API config
GOOGLE_MAPS_API_KEY=${google_key}
EOF
chown ec2-user:ec2-user /home/ec2-user/serve/server/.env
pushd /home/ec2-user/serve
systemctl start docker
docker compose -f docker-compose-dev.yml up -d
systemctl start nginx
systemctl start serve-be
