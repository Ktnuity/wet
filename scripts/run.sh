#!/bin/bash

set -e

"$(dirname "${BASH_SOURCE[0]}")/build.sh"

chmod +x "$(dirname "${BASH_SOURCE[0]}")/../bin/wet"

"$(dirname "${BASH_SOURCE[0]}")/../bin/wet" "$@"
