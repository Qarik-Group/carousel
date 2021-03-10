run:
	go run .

test:
	ginkgo watch ./...

gen:
	go generate ./...
