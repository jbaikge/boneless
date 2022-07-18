SRC := $(wildcard lambda/*)
LAMBDAS := $(notdir $(SRC))
ASSETS := $(addsuffix /handler,$(addprefix assets/,$(LAMBDAS)))

.PHONY: all deploy synth

all: $(ASSETS)

deploy: $(ASSETS)
	$(MAKE) -C _stack deploy

diff: $(ASSETS)
	$(MAKE) -C _stack diff

synth: $(ASSETS)
	$(MAKE) -C _stack synth

assets/%/handler: ./lambda/%/*.go *.go
	CGO_ENABLED=0 go build -o $@ $<
