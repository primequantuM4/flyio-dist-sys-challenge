#!/bin/bash
set -e

go install .
cd ~/maelstrom

./maelstrom test -w unique-ids --bin ~/go/bin/maelstrom-unique-ids \
    --time-limit 30 --rate 1000 --node-count 3 \
    --availability total --nemesis partition
