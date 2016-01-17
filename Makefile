clean:
	rm pebble-shutdown

build:
	GOOS=linux GOARCH=arm go build

.PHONY: clean build
