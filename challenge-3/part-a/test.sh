#!/bin/bash
set -e

go install .
cd ~/maelstrom

./maelstrom test -w broadcast --bin ~/go/bin/challenge-3-part-a \ 
--node-count 1 --time-limit 20 --rate 10
