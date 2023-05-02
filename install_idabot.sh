#!/bin/bash

# Variables
GO_PROGRAM_REPO="github.com/mfmayer/idabot" # Repository of the Go program
REPO_NAME="$(basename ${GO_PROGRAM_REPO})"
USERNAME="${REPO_NAME}_user"
GO_PROGRAM_NAME="${REPO_NAME}"
SERVICE_NAME="${REPO_NAME}"

# Find the absolute path of the 'go' binary
GO_BINARY_PATH=$(bash -l -c "which go")

if [ -z "${GO_BINARY_PATH}" ]; then
    echo "Go is not installed or not found in PATH. Please install Go and try again."
    exit 1
fi

# Prompt for environment variables
read -p "Enter DISCORD_BOT_TOKEN: " DISCORD_BOT_TOKEN
read -p "Enter OPENAI_API_KEY: " OPENAI_API_KEY
read -p "Enter AUTHORIZED_CHAT_PARTNER_ID: " AUTHORIZED_CHAT_PARTNER_ID

# Save original GOPATH
ORIGINAL_GOPATH="${GOPATH}"

# Set new GOPATH
export GOPATH="/home/${USERNAME}/go"

# Create new user
sudo useradd -m -s /bin/bash ${USERNAME}

# Install Go program using the absolute path of the 'go' binary
sudo -u ${USERNAME} ${GO_BINARY_PATH} install ${GO_PROGRAM_REPO}@latest

# Create systemd service file
sudo bash -c "cat > /etc/systemd/system/${SERVICE_NAME}.service << EOL
[Unit]
Description=My Go Service
After=network.target

[Service]
User=${USERNAME}
Group=${USERNAME}
Environment="GOPATH=/home/${USERNAME}/go"
Environment="DISCORD_BOT_TOKEN=${DISCORD_BOT_TOKEN}"
Environment="OPENAI_API_KEY=${OPENAI_API_KEY}"
Environment="AUTHORIZED_CHAT_PARTNER_ID=${AUTHORIZED_CHAT_PARTNER_ID}"
ExecStart=/home/${USERNAME}/go/bin/${GO_PROGRAM_NAME}
Restart=always

[Install]
WantedBy=multi-user.target
EOL"

# Set systemd service file permissions (root read-only)
sudo chmod 600 /etc/systemd/system/${SERVICE_NAME}.service

# Reload systemd daemon
sudo systemctl daemon-reload

# Enable and start service
sudo systemctl enable ${SERVICE_NAME}
sudo systemctl start ${SERVICE_NAME}

# Show service status
sudo systemctl status ${SERVICE_NAME}

# Reset GOPATH
export GOPATH="${ORIGINAL_GOPATH}"
