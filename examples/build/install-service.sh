#!/bin/bash

# Define the path to the binary and the service name
DEFAULT_SERVICE_NAME="publisher"
DEFAULT_SERVICE_DESCRIPTION="A Publisher"

# Allow specifying publisher binary name and/or description on the command line
# example:
# install-service.sh install my-publisher "My publisher description"
SERVICE_NAME="${2:-$DEFAULT_SERVICE_NAME}"
SERVICE_DESCRIPTION="${3:-$DEFAULT_SERVICE_DESCRIPTION}"

SERVICE_USER="root"
SERVICE_GROUP="root"
BINARY_WORKING_DIRECTORY="/opt/$SERVICE_NAME"
BINARY_PATH="/usr/bin/$SERVICE_NAME"

# Define the systemd service file content
SERVICE_FILE="[Unit]
Description=$SERVICE_DESCRIPTION
After=network.target

[Service]
User=$SERVICE_USER
Group=$SERVICE_GROUP
Type=simple
WorkingDirectory=$BINARY_WORKING_DIRECTORY
ExecStart=$BINARY_PATH
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target"

# Function to install the systemd service
function install_service() {
    # Check if the systemd service file already exists
    if [ -f /etc/systemd/system/$SERVICE_NAME.service ]; then
        echo "Systemd service file already exists. Exiting..."
        exit 1
    fi

    # Create the systemd service file
    echo "$SERVICE_FILE" > /etc/systemd/system/$SERVICE_NAME.service

    # Reload the systemd daemon to pick up the new service file
    systemctl daemon-reload

    # Enable the service to start on boot
    systemctl enable $SERVICE_NAME

    # Start the service
    systemctl start $SERVICE_NAME

    echo "Service installed successfully."
}

# Function to remove the systemd service
function remove_service() {
    # Check if the systemd service file exists
    if [ ! -f /etc/systemd/system/$SERVICE_NAME.service ]; then
        echo "Systemd service file does not exist. Exiting..."
        exit 1
    fi

    # Stop the service and disable it from starting on boot
    systemctl stop $SERVICE_NAME
    systemctl disable $SERVICE_NAME

    # Remove the systemd service file
    rm /etc/systemd/system/$SERVICE_NAME.service

    # Reload the systemd daemon to pick up the changes
    systemctl daemon-reload

    echo "Service removed successfully."
}

# Check if the user has provided an argument
if [ $# -ne 1 ]; then
    echo "Usage: $0 [install|remove]"
    exit 1
fi

# Parse the user's argument
case "$1" in
    install)
        install_service
        ;;
    remove)
        remove_service
        ;;
    *)
        echo "Usage: $0 [install|remove]"
        exit 1
esac
