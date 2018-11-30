.PHONY: all clean test cover release over-html

release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@git flow release finish $(version) -p -m "Release v$(version)"
	@git checkout develop
	@echo "Release v$(version) finished."

all: coverage.out

coverage.out: $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@CGO_ENABLED=0 go test -cover -coverprofile ./coverage.out.tmp ./...
	@cat ./coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > ./coverage.out
	@rm ./coverage.out.tmp

test: coverage.out

cover: coverage.out
	@echo ""
	@go tool cover -func ./coverage.out

cover-html: coverage.out
	@go tool cover -html=./coverage.out

clean:
	@rm ./coverage.out
	@go clean -i ./...

generate:
	@CGO_ENABLED=0 go generate ./...

lint:
	@golangci-lint run ./...
