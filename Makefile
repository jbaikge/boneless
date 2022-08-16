# SRC := $(wildcard lambda/*)
# LAMBDAS := $(notdir $(SRC))
# ASSETS := $(addsuffix /handler,$(addprefix assets/,$(LAMBDAS)))
# ADMIN := _frontend-admin/build/index.html

# .PHONY: all deploy diff frontend synth deploy-local diff-local synth-local localstack

# all: $(ASSETS) $(ADMIN)

# deploy: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack deploy

# deploy-local: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack $@

# diff: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack diff

# diff-local: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack $@

# frontend:
# 	$(MAKE) -C _frontend-admin start

# localstack:
# 	/usr/bin/docker run --rm --add-host host.docker.internal:host-gateway --publish 4566:4566 --publish 4510-4559:4510-4559 --env DEBUG=1 localstack/localstack:latest

# synth: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack synth

# synth-local: $(ASSETS) $(ADMIN)
# 	$(MAKE) -C _stack $@

# assets/%/handler: ./lambda/%/*.go *.go
# 	CGO_ENABLED=0 go build -o $@ $<

# _frontend-admin/build/index.html: _frontend-admin/src/*
# 	$(MAKE) -C _frontend-admin

LAMBDAS := $(wildcard lambda/*)
HANDLERS := $(addsuffix /handler,$(LAMBDAS))

.PHONY: all admin clean lambdas

all: lambdas admin

admin:
	$(MAKE) -C _frontend-admin

clean:
	rm --force $(HANDLERS)
	$(MAKE) -C _frontend-admin $@
	$(MAKE) -C _stack $@

deploy: all
	$(MAKE) -C _stack $@

deploy-admin: admin
	$(MAKE) -C _stack $@

deploy-repository: lambdas
	$(MAKE) -C _stack $@

lambdas: $(HANDLERS)

lambda/%/handler: lambda/%/*.go *.go
	CGO_ENABLED=0 go build -o $@ ./$(dir $<)
