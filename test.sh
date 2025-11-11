#!/bin/bash

set -e

./build.sh

echo "[=] Type Checker"
./wet ./demo/typecheck/init.wet

echo "[=] Arithmetics"
./wet ./demo/arithmetics/init.wet

echo "[=] Lang Test"
./wet ./demo/langtest/init.wet

echo "[=] Procedures"
./wet ./demo/proc/init.wet

echo "[=] Resource"
./wet ./demo/resource/init.wet

echo "[=] Tests"
./wet ./demo/tests/std.wet

echo "[=] Tokens"
./wet ./demo/tokens/init.wet

echo "[=] Tools"
./wet ./demo/tools/init.wet

echo "[=] Until + Unless"
./wet ./demo/untilunless/init.wet

echo "[=] Zip"
./wet ./demo/zip/init.wet

echo "[=] Commands"
echo "[help]"
./wet --help
echo "[license]"
./wet --license
echo "[version]"
./wet --version

echo "[=] No Args"
./wet

echo "[0] Success"
