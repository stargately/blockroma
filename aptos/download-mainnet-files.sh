#!/bin/bash

# Download mainnet genesis and waypoint files
echo "Downloading mainnet genesis blob..."
curl -o config/genesis.blob https://raw.githubusercontent.com/aptos-labs/aptos-networks/main/mainnet/genesis.blob

echo "Downloading mainnet waypoint..."
curl -o config/waypoint.txt https://raw.githubusercontent.com/aptos-labs/aptos-networks/main/mainnet/waypoint.txt

echo "Files downloaded successfully!"
ls -lh config/genesis.blob config/waypoint.txt