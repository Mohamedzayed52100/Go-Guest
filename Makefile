SHELL := /bin/bash

.PHONY: gen*

# detect if the CI mode is enabled
IS_CI_ENABLED := $(shell [ -z "$$CI" ] && v=false || v=true; echo $$v)

OS := $(shell uname)
OS_ALT := $(shell if [ "$(OS)" == "Darwin" ]; then echo "osx"; else echo $(OS); fi )
ARCH := $(shell uname -m)
ARCH_ALT := $(shell if [ "$(ARCH)" == "arm64" ]; then echo "aarch_64"; else echo $(ARCH); fi )

## START API/Proto generator targets ##
gen:
	$(MAKE) gen-grpc-stubs

gen-grpc-stubs:
	@for f in $$(find api -name "*.proto"); do \
		file_base=$${f%.proto}; \
		if [ ! -f $$file_base.pb.go -o $$f -nt $$file_base.pb.go ]; then \
			$$(protoc \
				-I $$HOME/.local/include \
				-I ./api \
				--go_out=./api \
				--go_opt=paths=source_relative \
				--go-grpc_out=require_unimplemented_servers=false:./api \
				--go-grpc_opt=paths=source_relative $$f);\
		fi \
	done

clean:
	find ./api -type f -name '*.mock.go' | xargs rm -f
	find ./api -type f -name '*.pb.go' | xargs rm -f
## END API/Proto generator targets ##

## START Migration generator targets ##
create-migration:
	if [ ! -x ./scripts/create_migration.sh ]; then chmod +x ./scripts/create_migration.sh; fi
	./scripts/create_migration.sh $(name) $(database)
## END Migration generator targets ##