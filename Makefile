SOURCES := $(shell find . -name "*.go")

all: build


.PHONY: build
build: dockerpatch


dockerpatch: $(SOURCES)
	go get ./...
	go build -o $@ ./cmd/$@


.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=6042 -workDir="$(realpath .)" -depth=-1
