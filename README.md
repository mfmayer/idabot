# idabot

Discord Chatbot with OpenAI API Integration

IDABot is a Discord chatbot written in Go that connects to the OpenAI API using the Go module github.com/user/repo. The chatbot requires the following environment variables to run: `DISCORD_BOT_TOKEN` (Discord bot token), `OPENAI_API_KEY` (OpenAI API key), and `AUTHORIZED_CHAT_PARTNER_ID` (Discord user ID to which the chatbot will respond). The provided installation script `install_idabot.sh` will set up the chatbot as a systemd service named `idabot`.

## Prerequisites

* Go programming language installed
* A Discord bot token
* An OpenAI API key
* A Discord user ID for the authorized chat partner
* Root access to execute the installation script

## Installation

The script has only been tested with an Ubuntu Server 22.02 installation. 

> ⚠️ **Use at your own risk and first check and understand the isntallation script's source code.**

1. Download the installation script install_idabot.sh from the repository.
2. Set the script's execution permissions:
    ```bash
    chmod +x install_idabot.sh
    ```
3. Run the installation script with root privileges:
    ```bash
    sudo ./install_idabot.sh
    ```
    The script will prompt you to enter the required environment variables: `DISCORD_BOT_TOKEN`, `OPENAI_API_KEY`, and `AUTHORIZED_CHAT_PARTNER_ID`.

4. Follow the prompts to enter the environment variables.

5. The installation script will create a new user, install the Go program, and set up a systemd service named `idabot` for the chatbot.

   1. After the installation is complete, the script will display the status of the `idabot` service.

## Usage

IDABot will automatically start upon installation and will be enabled to run at system startup. You can check the status of the service with:

```bash
sudo systemctl status idabot
```

To stop the service, run:

```bash
sudo systemctl stop idabot
```

To disable the service from running at system startup, run:

```bash
sudo systemctl disable idabot
```

## Troubleshooting

If you encounter any issues, check the logs for the idabot systemd service using:

```bash
sudo journalctl -u idabot
```

If you need to update the environment variables, edit the systemd service file located at `/etc/systemd/system/idabot.service` and reload the systemd daemon with:

```bash
sudo systemctl daemon-reload
```

Then restart the service:

```bash
sudo systemctl restart idabot
```