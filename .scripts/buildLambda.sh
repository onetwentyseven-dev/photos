#!/bin/bash

set -e
# store currect dir in a variable
current_dir=$(pwd)
lambda_prefix="$current_dir/cmd/lambda"
bin_dir="$current_dir/bin"

rm -rf "$bin_dir"
mkdir -p "$bin_dir"

for lambda in "$lambda_prefix"/*; do
	if [ -d "$lambda" ]; then
		lambda_name=$(basename "$lambda")
        echo "Building $lambda_name"
		CGO_ENABLED=0 go build -o "$bin_dir/$lambda_name/bootstrap" "$lambda_prefix/$lambda_name"
		cd "$bin_dir/$lambda_name"
		zip -qr "$lambda_name.zip" bootstrap
		mv "$lambda_name.zip" "$bin_dir"
		cd ../
		rm -rf "$lambda_name"
	fi
done

aws-vault exec --no-session ots -- deploy-functions $bin_dir