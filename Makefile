clean:
	rm pebble-shutdown

build:
	GOOS=linux GOARCH=arm go build

test:
	go test -v

deploy: build
	scp pebble-shutdown justin@chip1:/home/justin/bin && \
	scp pebble-shutdown justin@chip2:/home/justin/bin

.PHONY: clean build test deploy
