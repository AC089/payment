DATE=$(shell date)
ID=$(shell id -u -n)
#VERSION=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
BRANCH=$(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')

LD_FLAGS=-ldflags="-X 'payment/pkg/version.Version=${VERSION}' -X 'payment/pkg/version.BuildUser=${ID}' -X 'payment/pkg/version.BuildDate=${DATE}' -X 'payment/pkg/version.Branch=${BRANCH}'"

.PHONY: build clean grpc lint test

all: clean build lint

build:
	@echo "make payment"
	go build -v ${LD_FLAGS} -o ./cmd/payment ./pkg

clean:
	@echo "clean payment"
	go clean ./pkg
	rm -f ./cmd/payment

grpc:
	@echo "generate protobuf grpc"
	# find pkg/router/grpc/pb -name *.proto | awk -F "/" '{print $$NF}' | xargs protoc --proto_path=pkg/router/grpc/pb --go_out=plugins=grpc:. payment.proto
	find pkg/router/grpc/pb -name *.proto | awk -F "/" '{print $$NF}' | xargs protoc  -I. -I pkg/router/grpc/pb --gofast_out=plugins=grpc,import_path=pb:pkg/router/grpc/pb

lint:
	@echo "go lint"
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml ./... || exit 1

test:
	@echo "unit test"
	# gotests -all -w ./pkg
	# go test ./pkg -v