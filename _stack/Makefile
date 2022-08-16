.PHONY: deploy deploy-admin deploy-repository diff synth deploy-local diff-local synth-local

# Live endpoints
diff:
	cdk diff

synth:
	cdk synth

# Deploy both the repository and admin frontend
deploy: deploy-repository deploy-admin

# Deploy just the admin frontend. If .env.production.local does not exist, then
# create it by deploying the repository
deploy-admin: ../_frontend-admin/.env.production.local
	$(MAKE) -C ../_frontend-admin
	cdk deploy AdminStack

# Deploy the repository. The outputs file will trigger a repository re-deploy
# if anything has actually changed in the stack or binaries
deploy-repository: cdk.out/RepositoryStack.outputs.json

# (Re-)create the outputs file if it does not exist or the stack or binaries
# have changed.
cdk.out/RepositoryStack.outputs.json: lib/* bin/* ../lambda/admin/handler
	mkdir -p cdk.out
	cdk --outputs-file $@ deploy RepositoryStack

# Create .env.production.local from the outputs of the repository deployment
../_frontend-admin/.env.production.local: cdk.out/RepositoryStack.outputs.json
	jq --raw-output '"REACT_APP_API_URL=" + .RepositoryStack.EndpointUrl | rtrimstr("/")' $< > $@

# Localstack endpoints
diff-local:
	cdklocal diff

synth-local:
	cdklocal synth

deploy-local:
	cdklocal deploy
