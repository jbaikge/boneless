.PHONY: deploy diff synth

all:
	@echo "Can't decide what to do for 'all' so you get to choose your adventure:"
	@echo "make clean"
	@echo "make diff"
	@echo "make synth"
	@echo "make deploy"

clean:
	rm --recursive --force cdk.out

diff:
	cdk diff

synth:
	cdk synth

deploy:
	cdk deploy --all

# Deploy both the repository and admin frontend
# deploy: deploy-repository deploy-admin

# Deploy just the admin frontend. If .env.production.local does not exist, then
# create it by deploying the repository
# deploy-admin: ../_frontend-admin/.env.production.local
# 	$(MAKE) -C ../_frontend-admin
# 	cdk deploy AdminStack

# Deploy the repository. The outputs file will trigger a repository re-deploy
# if anything has actually changed in the stack or binaries
# deploy-repository: cdk.out/RepositoryStack.outputs.json

# (Re-)create the outputs file if it does not exist or the stack or binaries
# have changed.
# cdk.out/RepositoryStack.outputs.json: lib/* bin/* ../lambda/admin/handler ../lambda/frontend/handler
# 	mkdir -p cdk.out
# 	cdk --outputs-file $@ deploy RepositoryStack

# Create .env.production.local from the outputs of the repository deployment
# ../_frontend-admin/.env.production.local: cdk.out/RepositoryStack.outputs.json
# 	jq --raw-output '"REACT_APP_API_URL=" + .RepositoryStack.EndpointUrl | rtrimstr("/")' $< > $@
