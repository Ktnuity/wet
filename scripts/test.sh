#!/bin/bash

set -e

echo "[=] Building"

"$(dirname "${BASH_SOURCE[0]}")/build.sh"

echo "[=] Testing"

"$(dirname "${BASH_SOURCE[0]}")/../bin/test"
