#!/bin/bash

# Set default flags
DO_RUN=false
DO_BUILD=false
DO_TEST=false

# Check for flags
for arg in "$@"; do
    case $arg in
        --build)
            DO_BUILD=true
            shift
            ;;
        --test)
            DO_TEST=true
            shift
            ;;
        --run)
            DO_RUN=true
            shift
            ;;
        --all)
            DO_BUILD=true
            DO_TEST=true
            DO_RUN=true
            ;;
        --)
            shift
            break
    esac
done

# Enable exit on error
set -e

## Building step
if [ "$DO_BUILD" = true ]; then
    ./scripts/build.sh
fi

## Testing step
if [ "$DO_TEST" = true ]; then
    ./scripts/test.sh
fi

## Running step
if [ "$DO_RUN" = true ]; then
    ./scripts/run.sh "$@"
fi
