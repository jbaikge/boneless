SRC := $(wildcard lambda/*)
LAMBDAS := $(notdir $(SRC))
ASSETS := $(addsuffix /handler,$(addprefix assets/,$(LAMBDAS)))

.PHONY: all deploy

all: $(ASSETS)

deploy: $(ASSETS)
	$(MAKE) -C stack deploy

assets/%/handler: ./lambda/%/*.go
	CGO_ENABLED=0 go build -o $@ $^
