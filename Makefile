clean:
	rm pebble-shutdown

build:
	GOOS=linux GOARCH=arm go build

test:
	go test -v

.PHONY: clean build test
