# profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

buildLambda:
	./.scripts/buildLambda.sh

buildEdge:
	./.scripts/buildEdge.sh

buildAll: buildLambda buildEdge
