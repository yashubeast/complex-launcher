.PHONY: all dev

all: dev

dev:
	echo -e "first\nsecond" | go run .
