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

deploy: lambdas admin
	$(MAKE) -C _stack $@

lambdas: $(HANDLERS)

lambda/%/handler: lambda/%/*.go *.go
	CGO_ENABLED=0 go build -o $@ ./$(dir $<)
