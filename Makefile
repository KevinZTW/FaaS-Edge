redeploy:
	faas-cli build -f ./api.yml
	faas-cli push -f ./api.yml
	faas-cli deploy -f ./api.yml