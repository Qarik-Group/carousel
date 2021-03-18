ifdef VERSION
docker_registry = starkandwayne/carousel-concourse:$(VERSION)
else
docker_registry = starkandwayne/carousel-concourse
endif

run:
	go run .

test:
	ginkgo watch ./...

gen:
	go generate ./...
docker:
	docker build -t $(docker_registry) .

publish: docker
	docker push $(docker_registry)

fmt:
	find . -name '*.go' | while read -r f; do \
		gofmt -w -s "$$f"; \
	done

.DEFAULT_GOAL := docker

.PHONY: go-mod docker-build docker-push docker test fmt
