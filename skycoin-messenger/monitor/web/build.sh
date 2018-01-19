#!/usr/bin/env bash

echo start build discovery
./build-discovery.sh
echo start build manager
./build-manager.sh
echo "done"
