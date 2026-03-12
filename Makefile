build:
	go build -o cage ./cmd/runctrl

run: build
	sudo ./cage run /bin/sh
