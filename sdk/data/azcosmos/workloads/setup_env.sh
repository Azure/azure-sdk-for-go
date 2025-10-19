#!/bin/bash
# cspell:disable
# setup_env.sh - Automates Azure Cosmos SDK scale testing environment setup
# Usage: bash setup_env.sh

set -e

# 1. System update and install dependencies
echo "[Step 1] System update and install dependencies: started."
sudo apt-get update
sudo apt  install golang-go
sudo apt install neovim
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
echo "[Step 1] System update and install dependencies: completed."