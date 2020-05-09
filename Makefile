main=smp.go

fmt:
	gofmt -w .

build:
	go build $(main)

install:
	go install $(main)