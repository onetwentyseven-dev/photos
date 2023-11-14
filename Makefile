# profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

buildUpload:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./bin/upload cmd/upload/*.go
	aws-vault exec --no-session ots -- aws s3 cp ./bin/upload s3://ddouglas-desktop/upload

buildLambda:
	./.scripts/buildLambda.sh

buildEdge:
	./.scripts/buildEdge.sh

buildAll: buildLambda buildEdge
