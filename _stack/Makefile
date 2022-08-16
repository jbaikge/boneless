.PHONY: deploy diff synth deploy-local diff-local synth-local

# Live endpoints
diff:
	cdk --profile jake diff

synth:
	cdk --profile jake synth

deploy:
	cdk --profile jake --outputs-file /tmp/repository-stack.json deploy RepositoryStack
	jq --raw-output '"REACT_APP_API_URL=" + .RepositoryStack.EndpointUrl | rtrimstr("/")' /tmp/repository-stack.json > ../_frontend-admin/.env.production.local
	$(MAKE) -C ../_frontend-admin
	cdk --profile jake deploy AdminStack

# Localstack endpoints
diff-local:
	cdklocal diff

synth-local:
	cdklocal synth

deploy-local:
	cdklocal deploy
