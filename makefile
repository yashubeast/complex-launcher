.PHONY: all dev 1k i

all: dev

dev:
	echo -e "first\nsecond" | go run .

1k:
	for i in $$(seq 1 1000); do head -c 16 /dev/urandom | base64 | head -c 16; echo; done | go run .

i:
	cp cl cl-toggle ~/.local/bin
	ls ~/.local/bin
