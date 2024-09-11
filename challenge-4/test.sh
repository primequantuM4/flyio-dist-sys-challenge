#!/bin/bash
set -e

go install .
cd ~/maelstrom
./maelstrom test -w g-counter --bin ~/go/bin/challenge-4 --node-count 3 --rate 100 --time-limit 20 --nemesis partition
