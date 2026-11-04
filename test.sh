#!/bin/bash

set -e

echo "[=] Building"

./build.sh

echo "[=] Testing"

./test
