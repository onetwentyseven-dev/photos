.SHELL := /bin/bash

profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

applya:
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform apply --auto-approve'

applyb:
	cd ../ && make buildEdge
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform apply'


applyba:
	cd ../ && make buildEdge
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform apply --auto-approve'



init:
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform init'

initu:
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform init -upgrade'

apply:
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform apply'

plan:
	/bin/bash -c 'source .env && aws-vault exec --no-session ${profile} -- terraform plan'
