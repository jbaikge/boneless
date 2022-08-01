SRC := $(wildcard lambda/*)
LAMBDAS := $(notdir $(SRC))
ASSETS := $(addsuffix /handler,$(addprefix assets/,$(LAMBDAS)))
ADMIN := _frontend-admin/build/index.html

.PHONY: all deploy diff frontend synth deploy-local diff-local synth-local

all: $(ASSETS) $(ADMIN)

deploy: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack deploy

deploy-local: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack $@

diff: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack diff

diff-local: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack $@

frontend:
	$(MAKE) -C _frontend-admin start

synth: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack synth

synth-local: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack $@

assets/%/handler: ./lambda/%/*.go *.go
	CGO_ENABLED=0 go build -o $@ $<

_frontend-admin/build/index.html: _frontend-admin/src/*
	$(MAKE) -C _frontend-admin
