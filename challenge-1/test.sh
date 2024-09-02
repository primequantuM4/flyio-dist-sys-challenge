#!/bin/bash
set -e

go install .
cd ~/maelstrom
./maelstrom test -w echo --bin ~/go/bin/challenge-1 --node-count 1 --time-limit 10
