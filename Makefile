.PHONY: dep
dep:
	go get -u github.com/golang/dep/... && \
	dep ensure

.PHONY: build
build:
	go build
