#!/bin/bash

lambda_prefix="./cmd/lambda"
bin_dir="./bin"

rm -rf "$bin_dir"
mkdir -p "$bin_dir"

for lambda in "$lambda_prefix"/*; do
	if [ -d "$lambda" ]; then
		lambda_name=$(basename "$lambda")
        echo "Building $lambda_name"
		CGO_ENABLED=0 go build -o "$bin_dir/$lambda_name/bootstrap" "$lambda_prefix/$lambda_name"
	fi
done