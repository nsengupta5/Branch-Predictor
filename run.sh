#!/bin/bash

# Trace file directory
trace_directory="trace_files"
pattern="*.out"

if [ ! -d "$trace_directory" ]; then
    exit 1
fi

for file in "$trace_directory"/$pattern; do
    if [ -f "$file" ]; then
        echo "Running $file with always_taken"
        ./branch_predictor "$file" configs/always_taken.json
        echo "Running $file with two_bit"
        ./branch_predictor "$file" configs/two_bit.json
        echo "Running $file with gshare"
        ./branch_predictor "$file" configs/gshare.json
        echo "Running $file with two_bit_profiled"
        ./branch_predictor "$file" configs/two_bit_profiled.json
    fi
done
