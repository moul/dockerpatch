.PHONY: convey
convey:
	go get github.com/smartystreets/goconvey
	goconvey -cover -port=6042 -workDir="$(realpath .)" -depth=-1
