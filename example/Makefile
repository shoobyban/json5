all:
	tinygo build -o main.wasi -target wasi .
	wasmtime main.wasi

clean:
	rm -f main.wasi

.PHONY: all clean