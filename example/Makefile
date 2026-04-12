GOTOOLCHAIN ?= go1.25.4
TINYGO = env GOTOOLCHAIN=$(GOTOOLCHAIN) tinygo

all: run

build:
	$(TINYGO) build -o main.wasi -target wasi .

run: build
	@if command -v wasmtime >/dev/null 2>&1; then \
		wasmtime main.wasi; \
	else \
		echo "Built main.wasi. Install wasmtime to run it."; \
	fi

clean:
	rm -f main.wasi

.PHONY: all build run clean