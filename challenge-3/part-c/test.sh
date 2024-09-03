#!/bin/bash
set -e

go install .
cd ~/maelstrom
./maelstrom test -w broadcast --bin ~/go/bin/challenge-3-part-c --node-count 5 --time-limit 20 --rate 10 --nemesis partition
