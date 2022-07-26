SRC := $(wildcard lambda/*)
LAMBDAS := $(notdir $(SRC))
ASSETS := $(addsuffix /handler,$(addprefix assets/,$(LAMBDAS)))
ADMIN := _frontend-admin/build/index.html

.PHONY: all deploy diff frontend synth

all: $(ASSETS) $(ADMIN)

deploy: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack deploy

diff: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack diff

frontend:
	$(MAKE) -C _frontend-admin start

synth: $(ASSETS) $(ADMIN)
	$(MAKE) -C _stack synth

assets/%/handler: ./lambda/%/*.go *.go
	CGO_ENABLED=0 go build -o $@ $<

_frontend-admin/build/index.html: _frontend-admin/src/*
	$(MAKE) -C _frontend-admin
