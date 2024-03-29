LAMBDAS := $(wildcard bin/lambda-*)
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

cmd/%/handler: cmd/%/*.go models/*.go services/*.go repositories/*/*.go
	CGO_ENABLED=0 go build -o $@ ./$(dir $<)
